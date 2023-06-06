pushd dappdemo
go-contract build . || exit 1
popd

pushd mixinproxy
go-contract build . || exit 1
popd

pushd mixinproxy/token
go-contract build . || exit 1
popd

pushd mtg.xin
go-contract build . || exit 1
popd

pushd dappdemo
go-contract build . || exit 1
popd

