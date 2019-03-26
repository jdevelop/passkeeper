package support

type ControlsStub struct {
	UpFunc   func()
	DownFunc func()
	SendFunc func()
}

func (f *ControlsStub) OnClickUp() {
	if f.UpFunc != nil {
		f.UpFunc()
	}
	return
}

func (f *ControlsStub) OnClickDown() {
	if f.DownFunc != nil {
		f.DownFunc()
	}
	return
}

func (f *ControlsStub) OnClickOk() {
	if f.SendFunc != nil {
		f.SendFunc()
	}
	return
}
