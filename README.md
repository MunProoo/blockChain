### BlockChain Core 개발에 대해 공부한 레포지토리입니다.

---
BlockChain의 Core 개발을 직접적으로 해보고, CLI환경에서 실행해봅니다.   
기술 스택 : Golang, MongoDB   

👉 다음과 같은 메인 기능이 있습니다.
1. CreateWallet
2. ConnectWallet
3. ChangeWallet
4. TransferCoin (개인 -> 개인)
5. MintCoin (관리자 -> 개인)

해싱은 `github.com/hacpy/go-ethereum` 라이브러리를 이용하였고,  
`머클트리 구조`를 사용하여 데이터의 무결성과 검증속도를 올렸습니다.

ref)   
[강의](https://www.inflearn.com/course/%EB%94%B0%EB%9D%BC%ED%95%98%EB%A9%B4%EC%84%9C-%EB%A7%8C%EB%93%9C%EB%8A%94-%EB%B8%94%EB%A1%9D%EC%B2%B4%EC%9D%B8-%EC%BD%94%EC%96%B4-golang)
