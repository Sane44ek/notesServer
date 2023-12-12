package list

type node struct {
	index int64 // уникальный индекс ноды. Необходим для того, чтобы можно было удалять ноды из списка
	data  interface{}
	next  *node
}
