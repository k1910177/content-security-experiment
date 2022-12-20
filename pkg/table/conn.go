package table

import conn "github.com/uecconsecexp/secexp2022/se_go/connector"

func (table *Table) ChugakuSend(client *conn.ChugakuClient) error {
	return client.SendTable(table.ToSlice())
}

func ChugakuReceive(client *conn.ChugakuClient) (*Table, error) {
	data, err := client.ReceiveTable()
	if err != nil {
		return nil, err
	}
	return New(data), nil
}

func (table *Table) YobikouSend(server *conn.YobikouServer) error {
	return server.SendTable(table.ToSlice())
}

func YobikouReceive(server *conn.YobikouServer) (*Table, error) {
	data, err := server.ReceiveTable()
	if err != nil {
		return nil, err
	}

	return New(data), nil
}
