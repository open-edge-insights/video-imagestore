HOST_SYSTEM = $(shell uname | cut -f 1 -d_)
SYSTEM ?= $(HOST_SYSTEM)
CXX = g++

# cert/key path for imagestore and CA
SERVER_CERTS = /etc/ssl/imagestore
CLIENT_CERTS = /etc/ssl/imagestore
CA_CERTS = /etc/ssl/ca

# ImageStore hostname/port
HOST = localhost
PORT = 50055
INPUT_FILE = test.png
OUTPUT_FILE = test_output.png

CPPFLAGS += `pkg-config --cflags protobuf grpc`
CXXFLAGS += -std=c++11
PROTOC = protoc
PROTOPATH = ../../protobuff
PROTOCPPPATH = $(PROTOPATH)/cpp
GRPC_CPP_PLUGIN = grpc_cpp_plugin
GRPC_CPP_PLUGIN_PATH ?= `which $(GRPC_CPP_PLUGIN)`
LDFLAGS += -L/usr/local/lib `pkg-config --libs protobuf grpc++ grpc`\
           -Wl,--no-as-needed -lgrpc++_reflection -Wl,--as-needed\
           -ldl
DEV_MODE = false

all: clientTest

clientTest:	clean system-check ImageStore.pb.o ImageStore.grpc.pb.o clientTest.o
			$(CXX) $(PROTOCPPPATH)/ImageStore.pb.o $(PROTOCPPPATH)/ImageStore.grpc.pb.o clientTest.o $(LDFLAGS) -o clientTest
			./clientTest $(HOST) $(PORT) $(SERVER_CERTS)/imagestore_client_certificate.pem \
			                             $(SERVER_CERTS)/imagestore_client_key.pem \
										 $(CA_CERTS)/ca_certificate.pem \
										 $(INPUT_FILE) \
										 $(OUTPUT_FILE) \
										 $(DEV_MODE)

ImageStore.pb.o:
	protoc -I $(PROTOPATH) --cpp_out=$(PROTOCPPPATH)/ $(PROTOPATH)/ImageStore.proto
	g++ $(CXXFLAGS) $(CPPFLAGS) -c -o $(PROTOCPPPATH)/ImageStore.pb.o $(PROTOCPPPATH)/ImageStore.pb.cc

ImageStore.grpc.pb.o:
	protoc -I $(PROTOPATH)/ --grpc_out=$(PROTOCPPPATH)/ --plugin=protoc-gen-grpc=$(GRPC_CPP_PLUGIN_PATH) $(PROTOPATH)/ImageStore.proto
	g++ $(CXXFLAGS) $(CPPFLAGS) -c -o $(PROTOCPPPATH)/ImageStore.grpc.pb.o $(PROTOCPPPATH)/ImageStore.grpc.pb.cc

clientTest.o:
	g++ $(CXXFLAGS) $(CPPFLAGS) -c -o clientTest.o clientTest.cc

clean:
	rm -f clientTest.o clientTest $(PROTOCPPPATH)/*.h $(PROTOCPPPATH)/*.o

# The following is to test your system and ensure a smoother experience.
# They are by no means necessary to actually compile a grpc-enabled software.

PROTOC_CMD = which $(PROTOC)
PROTOC_CHECK_CMD = $(PROTOC) --version | grep -q libprotoc.3
PLUGIN_CHECK_CMD = which $(GRPC_CPP_PLUGIN)
HAS_PROTOC = $(shell $(PROTOC_CMD) > /dev/null && echo true || echo false)
ifeq ($(HAS_PROTOC),true)
HAS_VALID_PROTOC = $(shell $(PROTOC_CHECK_CMD) 2> /dev/null && echo true || echo false)
endif
HAS_PLUGIN = $(shell $(PLUGIN_CHECK_CMD) > /dev/null && echo true || echo false)

SYSTEM_OK = false
ifeq ($(HAS_VALID_PROTOC),true)
ifeq ($(HAS_PLUGIN),true)
SYSTEM_OK = true
endif
endif

system-check:
ifneq ($(HAS_VALID_PROTOC),true)
	@echo " DEPENDENCY ERROR"
	@echo
	@echo "You don't have protoc 3.0.0 installed in your path."
	@echo "Please install Google protocol buffers 3.0.0 and its compiler."
	@echo "You can find it here:"
	@echo
	@echo "   https://github.com/google/protobuf/releases/tag/v3.0.0"
	@echo
	@echo "Here is what I get when trying to evaluate your version of protoc:"
	@echo
	-$(PROTOC) --version
	@echo
	@echo
endif
ifneq ($(HAS_PLUGIN),true)
	@echo " DEPENDENCY ERROR"
	@echo
	@echo "You don't have the grpc c++ protobuf plugin installed in your path."
	@echo "Please install grpc. You can find it here:"
	@echo
	@echo "   https://github.com/grpc/grpc"
	@echo
	@echo "Here is what I get when trying to detect if you have the plugin:"
	@echo
	-which $(GRPC_CPP_PLUGIN)
	@echo
	@echo
endif
ifneq ($(SYSTEM_OK),true)
	@false
endif