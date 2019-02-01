package ln

type wrappedError {
  err error
}

func (we wrappedError) Error() string {
  return ".i ko cikre ti"
}

func (we wrappedError) Unwrap() error {
  return we.err
}

func TestErrorsHack(t *testing.T) {
  f := F{}
  doGoError(errors.New("test"), f)
  
  if _, ok := f["wrapped_err"]; ok {
    t.Error("didn't expect f[wrapped_err] to be set")
  }
  
  f = F{}
  err := errors.New("test")
  doGoError(wrappedError{err: err}, f)
  
  if gerr, ok := f["wrapped_err"]; ok && gerr != err {
    t.Logf("f[wrapped_err]: %v", gerr)
    t.Logf("err:            %v", err)
    t.Errorf("expected f[wrapped_err] to return err, see -v")
  }
}