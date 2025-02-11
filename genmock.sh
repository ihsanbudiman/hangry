# make me a generate mock file using go mock gomock
mockgen -source=./usecase/*.go -destination=./mock/mock_uscase.go -package=mock