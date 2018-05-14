package basics




type MysqlRepository struct {

}

func (r *MysqlRepository)SelectById(id,pSlice interface{}) bool  {
	return false
}
func (r *MysqlRepository)Select(search interface{},pSlice interface{}) bool {
	return false
}
func (r *MysqlRepository)SelectMany(search interface{},sortKey string, pSlice interface{}) bool {
	return false
}

func (r *MysqlRepository)InsertOrUpdate(pSlice interface{}) bool {
	return false
}

func (r *MysqlRepository)Insert(pSlice interface{}) bool{
	return false
}
func (r *MysqlRepository)Update(id,pData interface{}) bool {
	return false
}