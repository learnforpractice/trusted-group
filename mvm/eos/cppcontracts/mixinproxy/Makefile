all : mixinproxy.wasm

mixinproxy.wasm : mixinproxy.cpp mixinproxy.hpp
	cdt-cpp mixinproxy.cpp

test : mixinproxy.wasm
# ipyeos -m pytest -x -s tests/test.py
# ipyeos -m pytest -x -s tests/test.py -k test_event
# ipyeos -m pytest -x -s tests/test.py -k test_expiration
# ipyeos -m pytest -x -s tests/test.py -k test_fee
	ipyeos -m pytest -x --log-level=1 --log-cli-level=ERROR -s tests/test.py
