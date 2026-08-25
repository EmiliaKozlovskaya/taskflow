#мультистейдж сборка - сборка нового docker образа на основе двух и более docker образов
#STAGE 1
FROM golang:1.25.0-bookworm AS builder
WORKDIR /app
#копирует файлы mod sum с компьютера в текущую директорию контейнера /app 
#COPY <откуда> <куда>
COPY go.mod go.sum ./
#смотрит в go.mod,go.sum и устанавливает необходимые оттуда зависимости
RUN go mod download 
#копируем все файлы проекта в рабочую директорию
COPY . . 
# go build -o(output) <путь_к_бинарнику> <путь_к_пакету>
# CGO_ENABLED=0 - чтобы собрать статический бинарник без зависимостей от C-библиотек, которые имеются в большом образе bookworm и не имеются в маленьком alpine, с помощью которого хотим запускать
RUN CGO_ENABLED=0 go build -o /app/cmd/todoapp/exe /app/cmd/todoapp

#STAGE 2
#для запуска обычно используют образ alpine
FROM alpine:3.24
WORKDIR /app
#нужно указать с какого stage хотим копировать
COPY --from=builder /app/cmd/todoapp/exe /app
COPY --from=builder /app/public /app/public
#для service.GetMainPage нужен project_root
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]
