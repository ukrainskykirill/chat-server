package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ChatsService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i MessagesService -o ./mocks/ -s "_minimock.go"
