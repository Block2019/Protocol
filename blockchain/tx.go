package blockchain

type TxOutput struct { //Transaction output
	Value  int    //Value of Transaction
	PubKey string //Wallet Address
}

type TxInput struct { //References to previous outputs
	ID        []byte //Hash Transaction ID
	Out       int    // index which the previous output appears
	Signature string // Script that provides the data that is sued in the PubKey
}

type Transaction struct {
	ID      []byte     //Hash transaction ID
	Inputs  []TxInput  //Values entered into transaction
	Outputs []TxOutput //Values shown by the transaction
}

//Function to unlock the data
func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data //Signature is the data that is used in the PubKey
}

//Function to unlock the data
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data //PubKey is the data that is used in the PubKey
}
