package test

//func doMstsc(address string, user string, pwd string) error {
//	var addressIp = getIp(address)
// 	if out, err := RunCmd(fmt.Sprintf("cmdkey /generic:%s /user:\"%s\" /pass:\"%s\"", addressIp, user, pwd)); err != nil
// {
//		fmt.Println("错误 " + err.Error())
//		return err
//	} else {
//		if !strings.Contains(out, "成功添加凭据") {
//			fmt.Println("错误 " + err.Error())
//			return errors.New(out)
//		}
//	}
//	defer func() { _, _ = RunCmd(fmt.Sprintf("cmdkey /delete:%s", addressIp)) }()
//	//mstsc /v:address
//	if _, err := RunCmd(fmt.Sprintf("mstsc /v:%s", address)); err != nil {
//		return err
//	}
//	return nil
//}
//func getIp(address string) string {
//	var ipIndex = strings.Index(address, ":")
//	addressIp := address
//	if ipIndex > 0 {
//		addressIp = address[:ipIndex]
//	}
//	return addressIp
//}
//func RunCmd(bash string) (string, error) {
//	//fmt.Println(bash)
//	cmd := exec.Command("cmd.exe")
//	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, bash), HideWindow: true}
//	out, err := cmd.Output()
//	if err != nil {
//		fmt.Println("错误", err)
//		return string(out), err
//	}
//	output, err := simplifiedchinese.GBK.NewDecoder().Bytes(out)
//	if err != nil {
//		fmt.Println("错误", err)
//		return string(output), err
//	}
//
//	return string(output), nil
//}
