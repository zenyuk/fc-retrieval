#!/bin/bash
#
# get-dependency-licenses.sh
# Finds GitHub dependencies recursively used in a Golang project,
# and creates a file summarising their licenses.
#
# David Hyland-Wood, 1 March 2021
# Copyright (c) 2021 ConsenSys Software Inc.

# Functions -------------------%<-----------------------

# Prints message to STDERR.
error() {
	echo -e "$@" 1>&2
}

# Prints message to STDERR and exits.
fail() {
	error "$@"
	exit 1
}

# Removes double quotes from string argument.
remove_quotes() {
	line=$@
	line=$(echo $line | sed s/\"//g)
	echo $line
}

# Prints messages to STDERR only if verbose mode requested.
ifverbose() {
	if [ ! -z "${verbose}" ]; then
		echo -e "$@" 1>&2
	fi
}

# Tries to determine the type of a license file.
check_license() {
	reentrant=$1
	license=$@
	output=""
	
	# Test the license against common titles.
	if [[ $license == *"Apache License"* ]]; then
		output=${output}"License type: Apache License"
	elif [[ $license == *"MIT License"* ]]; then
		output=${output}"License type: MIT License"
	elif [[ $license == *"BSD-style license"* ]]; then
		output=${output}"License type: BSD-style License"
	elif [[ $license == *"SOFTWARE IS PROVIDED"* && $license == *"AS IS"* ]]; then
		# NB: This check only works if it is after the MIT License check.
		output=${output}"License type: BSD-style License"
	elif [[ $license == *"GNU LESSER GENERAL PUBLIC LICENSE"* ]]; then
		output=${output}"License type: GNU LGPL License"
	elif [[ $license == *"GNU GENERAL PUBLIC LICENSE"* && $license == *"Version 3"* ]]; then
		output=${output}"License type: GNU GPL v3 License"
	elif [[ $license == *"GNU GENERAL PUBLIC LICENSE"* && $license == *"Version 2"* ]]; then
		# Add WARNING!
		error "\tWARNING: License problem identified: GPL v2. Check LEGAL.txt"
		output=${output}"License type: WARNING: GPL v2 license identified. RESOLVE BEFORE SHIPPING!\n"
		output=${output}"License type: GNU GPL v2 License"
	elif [[ $license == *"404: Not Found"* ]]; then
		if [[ $reentrant == *"checkagain"* ]]; then
			# Add WARNING!
			ifverbose "\tWARNING: Second time license check failed."
			error "\tWARNING: License problem identified: 404 Not Found. Check LEGAL.txt"
			output=${output}"WARNING: No license identified. RESOLVE BEFORE SHIPPING!\n"
			output=${output}"License type: 404 Not Found"
		else
			ifverbose "\tWARNING: First time license check failed."
			output=${output}"404 Not Found"
		fi
	else
		# Add WARNING!
		error "\tWARNING: License problem identified: Unknown License. Check LEGAL.txt"
		output=${output}"WARNING: No license identified. RESOLVE BEFORE SHIPPING!\n"
		output=${output}"License type: Unknown license! \n$license\n"
	fi
	echo $output
}

# END Functions -------------------%<-----------------------

# Capture the arguments array.
args=("$@")

# If no command line options are given, process all repositories.
if [ ${#args[@]} -eq 0 ]; then
    args+=("-a")
fi

# Determine whether we should be in verbose mode.
verbose=
if [[ " ${args[0]} " == *"-v"* ]]; then
    error "$0: Processing in verbose mode"
	verbose=1
fi

# List all GitHub repository short names in this project.
declare -a shortrepositories
shortrepositories+=("common")
shortrepositories+=("docs")
shortrepositories+=("gateway")
shortrepositories+=("gateway-admin")
shortrepositories+=("provider")
shortrepositories+=("provider-admin")
shortrepositories+=("register")

# This array will hold all GitHub repository names to process.
declare -a repositories

if [[ " ${args[0]} " == *"-u"* || " ${args[0]} " == *"-h"* ]]; then
    echo -e "Usage:"
	echo -e "\t$0 [-u|-h]"
	echo -e "\t\tPrint usage message and exit"
	echo -e "\t$0 -a"
	echo -e "\t\tProcess all repositories"
	echo -e "\t$0 [one or more short repository names]"
	echo -e "\t\tProcess only named repositories, e.g."
	echo -e "\t\t$0 common gateway # Process only the "
	echo -e "\t\tfc-retrieval-common and fc-retrieval-gateway repositories"
	echo
	echo -e "Available repository short names:"
	for shortname in "${shortrepositories[@]}"
	do
		echo -e "\t$shortname"
	done
	exit
elif [[ " ${args[0]} " == *"-a"* ]]; then
	error "Processing all repositories"
	for shortname in "${shortrepositories[@]}"
	do
		ifverbose "\t$shortname"
		repositories+=("fc-retrieval-$shortname")
	done
else
	error "Processing named repositories"
	for arg in "${args[@]}"
	do
		# Add the argument to the list of repositories to process if it looks legit.
		# Do this ONLY if $arg is present in "${shortrepositories[@]}"
		for repo in "${shortrepositories[@]}"
		do
			if [ "$repo" == "$arg" ]; then
			    repositories+=("fc-retrieval-$arg")
			fi
		done
	done
fi

# Report the repositories being processed.
echo -e "Processing these repositories:"
for repo in "${repositories[@]}"
do
	error "\t$repo"
done

# Ensure that we are operating at the base directory for the current repo
cd `dirname "$0"` # Go to the directory this script is located in.
cd ../.. # Go two levels up; this should be the directory containing all
         # the fc-retrieval-* repository checkouts.

for localcheckout in "${repositories[@]}"
do
	if [ -d "$localcheckout" ]; then

		error "---------------------------------------------"
		error "Processing $localcheckout"
		error "\tlocal checkout found."
		
		# Process the currrent checkout.
		cd $localcheckout
		
		# Find recursively all the GitHub URLs used as dependencies, and write
		# them to a temporary file. Avoid any alternative references to github.com.
		grep -r -h 'github.com' `ls -d */` | grep -v echo | grep -v 'go get' | grep -v '^\s*//' | grep -v '^\s*#' | grep -v 'Ignoring' | grep -v 'projecturl' | sed 's/\t+//g' | sed 's/\s+//g' | sed 's/^.*"github/github/'| sed 's/\s*$//' | sed 's/\r$//' | sed 's/"$//' | sort | uniq > .dependencyurls.txt

		# Create an array to hold the license URLs
		declare -a license_urls

		# Create an array to hold the project URLs
		declare -a project_urls

		# Process the temporary file line by line to create license URLs from
		# the dependency URLs.
		while read depurl; do
	
			ifverbose "Processing URL: $depurl"
	
			# Split the URL into an array based on the '/' separator.
			IFS='/' read -ra url_components <<< "$depurl"
	
			# Test the URL components to avoid errors
			if [[ -z "${url_components[1]}" || "${url_components[1]}" == *".."* ]]; then
				ifverbose "\tIgnoring github.com reference: var for org is blank in $depurl";
				continue
			fi
			if [[ -z "${url_components[2]}" || "${url_components[2]}" == *".."* ]]; then
				ifverbose "\tIgnoring github.com reference: var for project is blank in $depurl";
				continue
			fi
	
			# Construct license URL
			licenseurl='"https://raw.githubusercontent.com/'${url_components[1]}'/'${url_components[2]}'/master/LICENSE"' 
	
			# Construct project URL
			projecturl='"https://github.com/'${url_components[1]}'/'${url_components[2]}'/"' 
	
			# De-duplicate the license URLs
			if [[ " ${license_urls[*]} " == *"$licenseurl"* ]];
			then
				ifverbose "\tStopping addition of duplicate license URL: $licenseurl"
			else
			    # Add license url to license array
				ifverbose "\tAdding URL to license array: $licenseurl"
				license_urls+=($licenseurl)
			fi
	
			# De-duplicate the project URLs
			if [[ " ${project_urls[*]} " == *"$projecturl"* ]];
			then
				ifverbose "\tStopping addition of duplicate project URL: $projecturl"
			else
			    # Add project url to project array
				ifverbose "\tAdding URL to project array: $projecturl"
				project_urls+=($projecturl)
			fi
	
		done <.dependencyurls.txt

		# Get the license for each license URL

		# Output to the file LEGAL.txt in the current repo from here.

		# Output the file header.
		rm LEGAL.txt
		echo "LEGAL.txt" >>LEGAL.txt
		echo "Copyright (c) 2021 ConsenSys Software Inc." >>LEGAL.txt
		echo "" >>LEGAL.txt
		echo "Project: $localcheckout" >>LEGAL.txt
		echo "" >>LEGAL.txt
		echo "This file lists dependencies used by this project, and provides" >>LEGAL.txt
		echo "attribution to open source software projects." >>LEGAL.txt
		echo "---------------------------------------------------------------" >>LEGAL.txt
		echo >>LEGAL.txt

		counter=0
		for licenseurl in "${license_urls[@]}"
		do
	
			# Display the project URL.
			projurl=${project_urls[$counter]}
			projectdisplay=$(remove_quotes $projurl)
			echo "Project: ${projectdisplay}" >>LEGAL.txt
			ifverbose "Project: ${projectdisplay}"
	
			# Retrieve the license file.
			cmd="curl -s -S $licenseurl"
			ifverbose "\t\tCommand: "$cmd
			license=$( eval $cmd )
			
			licensetype=$( check_license $license )
			if [[ $licensetype == *"404 Not Found"* ]]; then
				# Check for LICENCE (LICEN*C*E) also to account for spelling differences.
				newlicenseurl=$licenseurl
				newurl=${newlicenseurl/LICENSE/LICENCE}
				echo "License URL: $newurl" >>LEGAL.txt
				ifverbose "\tEncountered 404. Checking $newurl as well."
				newcmd="curl -s -S $newurl"
				newlicense=$( eval $newcmd )
				licensetype=$( check_license "checkagain" $newlicense )
			else
				# Display the license URL.
				licensedisplay=$(remove_quotes $licenseurl)
				echo "License URL: $licensedisplay" >>LEGAL.txt
				ifverbose "License: ${licensedisplay}"
			fi
			echo -e $licensetype >>LEGAL.txt
			
			echo >>LEGAL.txt # Add a blank line between entries.
			counter=$((counter+1))
		done

		# Remove temporary file for this repository
		rm .dependencyurls.txt
		
		# Return to the base directory level.
		cd ..
		
		# Finished processing the current checkout.
		error "\tLEGAL.txt file written."
		error "---------------------------------------------"
		
	else
		error "---------------------------------------------"
		error "WARNING: $localcheckout not found"
		error "\t... skipping"
		error "---------------------------------------------"
		error
	fi
done 
exit # Exit cleanly
