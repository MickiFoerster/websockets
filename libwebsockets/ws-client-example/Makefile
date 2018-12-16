OPENSSL_INCLUDE=-I $$OPENSSL_BUILD/include 
OPENSSL_LIB=$$OPENSSL_BUILD/lib/libssl.a  $$OPENSSL_BUILD/lib/libcrypto.a 
LIBWEBSOCKET_INCLUDE=-I$$LIBWEBSOCKET_BUILD/include
LIBWEBSOCKET_STATIC=$$LIBWEBSOCKET_BUILD/lib/libwebsockets.a

minimal-ws-client-echo : minimal-ws-client-echo.c
	clang -o $@ \
		    $(OPENSSL_INCLUDE) \
				$(LIBWEBSOCKET_INCLUDE) \
				$^ \
				$(LIBWEBSOCKET_STATIC) \
			  $(OPENSSL_LIB) \
				-pthread -ldl
