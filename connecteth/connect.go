package connecteth

import(
    "fmt"
    "context"
	"encoding/json"
    "math/big"
    "time"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/rpc"
)

type Client struct {

    RpcClient *rpc.Client
    EthClient *ethclient.Client
}

func Connect(host string) (*Client, error) {

    rpcClient, err := rpc.Dial(host)

    if err != nil {

        return nil, err
    }

    ethClient := ethclient.NewClient(rpcClient)

    return &Client{rpcClient, ethClient}, nil
}

func (ec *Client) GetBlockNumber(ctx context.Context) (*big.Int, error) {

    var result hexutil.Big
    err := ec.RpcClient.CallContext(ctx, &result, "eth_blockNumber")
    return (*big.Int)(&result), err
}

type Message struct {
    To *common.Address `json:"to"`
    From *common.Address `json:"from"`
    Value string `json:"value"`
    GasLimit string `json:"gas"`
    GasPrice string `json:"gasPrice"`
    Data []byte `json:"data"`
}

func (msg *Message) String() string {
    if str, err := json.Marshal(msg); err != nil {
        panic(err)
    }else {
        return string(str)
    }
}

func NewMessage(from *common.Address, to *common.Address, value *big.Int, gasLimit *big.Int, gasPrice *big.Int, data []byte) Message {
	return Message{
		From: from,
		To: to,
		Value: toHexInt(value),
		GasLimit: toHexInt(gasLimit),
		GasPrice: toHexInt(gasPrice),
		Data: data,
	}
}

func (ec *Client) SendTransaction(ctx context.Context, tx *Message) (common.Hash, error) {
    var txHash common.Hash
    err := ec.RpcClient.CallContext(ctx, &txHash, "eth_sendTransaction")
    return txHash, err
}

func (ec *Client) CheckTransaction(ctx context.Context, receiptChan chan *types.Receipt, txHash common.Hash, retrySeconds time.Duration) {
    // check transaction receipt
    go func() {
        fmt.Printf("Check transaction: %s\n", txHash.String())
        for {
            receipt, _ := ec.EthClient.TransactionReceipt(ctx, txHash)
            if receipt != nil {
                receiptChan <- receipt
                break
            }else {
                fmt.Printf("Retry after %d second\n", retrySeconds)
                time.Sleep(retrySeconds *time.Second)
            }
        }
    }()
}

func toHexInt(n *big.Int) string {
    return fmt.Sprintf("%x", n)
}

