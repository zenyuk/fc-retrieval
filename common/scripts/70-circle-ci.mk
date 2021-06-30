#
# make ci: it will run steps, which runs on circle CI before normal itest
#  

ci-1 : deps build-servers micro-images coverage

ci-2 : lotus-images

ci-3 : build-servers-test

ci-4 : itest-poc1

# a quick check
# skip lotus and any itest
ci: ci-1 ci-3

# ci more
cimore: ci-1 ci-2 ci-3 ci-4

