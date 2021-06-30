
default:

ci CI: deps build-servers micro-images coverage

-include ./common/scripts/*.mk
-include ../local*.mk
