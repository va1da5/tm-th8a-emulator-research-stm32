.PHONY: help
help: 
	@echo "make wheel - build and flash emulator node"
	@echo "make shifter - build and flash shifter node"
	@echo "make help - show this message"


.PHONY: wheel
wheel:
	cd ./wheel && make all


.PHONY: shifter
shifter:
	cd ./shifter && make all
