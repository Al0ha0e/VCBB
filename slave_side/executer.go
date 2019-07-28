package slave_side

type Executer interface {
	Run()
}

type PyExecuter struct {
	code string
}

func (this *PyExecuter) Run() {

}
