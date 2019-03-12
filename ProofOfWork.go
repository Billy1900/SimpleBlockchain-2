package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const TargetBit = 16//此处指hash值前16位为0，此处将难度设置为全局变量

//pow struct
type ProofOfWork struct {
	Block *Block//验证的blocks
	Target *big.Int//target hash
}

//创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target,256-TargetBit)
	return &ProofOfWork{block,target}
}

//返回hash与nonce
func (pow *ProofOfWork) Run() ([]byte,int64)  {
	nonce := 0
	hashInt := new(big.Int)
	var hash [32]byte
	for{
		dataBytes := pow.prepareData(nonce)//字符拼接
		hash = sha256.Sum256(dataBytes)//计算hash
		fmt.Printf("\r%d:%x",nonce,hash)
		hashInt.SetBytes(hash[:])
		/* Cmp
		compare x and y:
		-1 if x < y
		0 if x==y
		1 if x > y
		*/
		if pow.Target.Cmp(hashInt) == 1 {
			break
		}
		nonce++
	}
	fmt.Println()
	return hash[:],int64(nonce)
}

func (pow *ProofOfWork) prepareData(nonce int) []byte{
	data := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.Data,
		IntToHex(pow.Block.TimeStamp),
		IntToHex(int64(TargetBit)),
		IntToHex(int64(nonce)),
	},[]byte{})
	return data
}

func (pow *ProofOfWork) IsValid() bool{
	hashInt := new(big.Int)
	hashInt.SetBytes(pow.Block.Hash)
	return pow.Target.Cmp(hashInt) == 1
}