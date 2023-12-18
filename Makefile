BUILD_DIRECTORY=build/
EXECUTABLE_FILENAME=api-userbase
EXECUTABLE_PATH=${BUILD_DIRECTORY}${EXECUTABLE_FILENAME}
PURPLE=\033[0;35m
NC=\033[0m

clean:
	@echo "${PURPLE}Cleaning build directory...${NC}"
	rm -fr ${BUILD_DIRECTORY}

install:
	@echo "${PURPLE}Updating modules...${NC}"
	cd src && go mod tidy

build: clean install
	@echo "${PURPLE}Building executable...${NC}"
	cd src && go build -o ../${EXECUTABLE_PATH} .

run: build
	@echo "${PURPLE}Running executable...${NC}"
	./${EXECUTABLE_PATH}
