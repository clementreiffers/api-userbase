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

build_alpine: clean install
	@echo "${PURPLE}Building executable for alpine...${NC}"
	cd src && GOOS=linux GOARCH=amd64 go build -o ../${EXECUTABLE_PATH} .
	@echo "${PURPLE}Building done!${NC}"

run: build
	@echo "${PURPLE}Running executable...${NC}"
	./${EXECUTABLE_PATH}

curls-create-user:
	curl -X POST -H 'Content-Type: application/json' -d '{"name":"clement", "ID": "srlkghfq"}' http://localhost:8080/add-user
	curl -X GET -H 'Content-Type: application/json' -d '{"name":"clement", "ID": "111"}' http://localhost:8080/get-user

curls: curls-create-user