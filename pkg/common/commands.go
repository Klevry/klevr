package common

func init() {
	InitCommand(testCommand())
}

// 샘플 커맨드 작성
func testCommand() Command {
	return Command{
		// 커맨드 명칭
		Name: "SampleCommand",
		// 커맨드 실행 로직 구현
		Run: func(param *map[string]interface{}) (interface{}, error) {
			return nil, nil
		},
		// 커맨드 중지 시 복구 로직 구현
		Recover: func(param *map[string]interface{}) (interface{}, error) {
			return nil, nil
		},
	}
}
