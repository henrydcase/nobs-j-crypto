package mkem

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/henrydcase/nobs/dh/sidh"
	"github.com/henrydcase/nobs/dh/sidh/common"
	"github.com/henrydcase/nobs/drbg"
)

var rng *drbg.CtrDrbg

func init() {
	var tmp [32]byte

	rand.Read(tmp[:])
	rng = drbg.NewCtrDrbg()
	if !rng.Init(tmp[:], nil) {
		panic("Can't initialize DRBG")
	}
}

type sikeVec struct {
	id   uint8
	name string
	kem  *KEM
	PkB  string
	PrB  string
	mkem MultiKEM
}

func init() {
	var tmp [32]byte

	rand.Read(tmp[:])
	rng = drbg.NewCtrDrbg()
	if !rng.Init(tmp[:], nil) {
		panic("Can't initialize DRBG")
	}
}

// testPanic returns true if call to function 'f' caused panic.
func testPanic(f func()) error {
	var hasPanicked = errors.New("no panic detected")
	defer func() {
		if r := recover(); r != nil {
			hasPanicked = nil
		}
	}()
	f()
	return hasPanicked
}

// helper
func IsOk(t testing.TB, f error, msg string) {
	t.Helper()
	if f != nil {
		t.Error(msg)
	}
}

var tdataSike = map[uint8]sikeVec{
	sidh.Fp434: {
		sidh.Fp434, "P-434", NewSike434(rand.Reader),
		"1BD0A2E81307B6F96461317DDF535ACC0E59C742627BAE60D27605E10FAF722D" +
			"22A73E184CB572A12E79DCD58C6B54FB01442114CBE9010B6CAEC25D04C16C5E" +
			"42540C1524C545B8C67614ED4183C9FA5BD0BE45A7F89FBC770EE8E7E5E391C7" +
			"EE6F35F74C29E6D9E35B1663DA01E48E9DEB2347512D366FDE505161677055E3" +
			"EF23054D276E817E2C57025DA1C10D2461F68617F2D11256EEE4E2D7DBDF6C8E" +
			"34F3A0FD00C625428CB41857002159DAB94267ABE42D630C6AAA91AF837C7A67" +
			"40754EA6634C45454C51B0BB4D44C3CCCCE4B32C00901CF69C008D013348379B" +
			"2F9837F428A01B6173584691F2A6F3A3C4CF487D20D261B36C8CDB1BC158E2A5" +
			"162A9DA4F7A97AA0879B9897E2B6891B672201F9AEFBF799C27B2587120AC586" +
			"A511360926FB7DA8EBF5CB5272F396AE06608422BE9792E2CE9BEF21BF55B7EF" +
			"F8DC7EC8C99910D3F800",
		"4B622DE1350119C45A9F2E2EF3DC5DF56A27FCDFCDDAF58CD69B903752D68C20" +
			"0934E160B234E49EDE247601",
		MultiKEM{}},
	sidh.Fp503: {
		sidh.Fp503, "P-503", NewSike503(rand.Reader),
		"4032A90B6C036B7D2A83878AD116641AD319E420235A505F3F5C3DEC27C87A6C" +
			"BA0792201D6E7B196C582D43CAF86CB2C7DEFA6598B543C946CDDF62EF9A328C" +
			"8719B66BA5052231DAE13AF7D9CDEBB4ED327773C7AE0818F41AF1D28CD78B16" +
			"C996232528235C8392B8FCFD925CB311B2A801B0402A90E527261EA32F2BEF67" +
			"7C544908D5509B8AB7D7BF20456727AD358AD585306A0B28F6B2AA583CE8A3E0" +
			"BB92D8CD55347D39D4E3C30D3D0F96EABB721A6968CDD143FE9227643CF697FB" +
			"2DF0B71322B5EA1505D0DDBF70A2FD1193011F3BC18AA1E127C614B76969DCDA" +
			"45A2072B519A1074FDA49F5C828450C6A007BF8D7CDDD5D2FC112119C679CA3A" +
			"B16C6960B25F6C681A7DCED0F0E3901740D3DBF3A33011EB7DA460E8ADA80EE3" +
			"45B2B71420950A9A803E4F11330EB91CCABB1EEE4D875A109D7724ABD201272C" +
			"0B4981BDCDFA70F3430A89D2A88EEED474CF0CFAC65CE883F44B4722FA280C6F" +
			"A9C4724D414B35AF69D6ECB21BFDA23BFF6B66C22C2451DC8E1C",
		"7BF6938C975658AEB8B4D37CFFBDE25D97E561F36C219A0E8FE645816DBBC7ED7B57" +
			"7700AE8DC3138E97A0C3F6F002065C92A0B1B8180208", MultiKEM{}},
	sidh.Fp751: {
		sidh.Fp751, "P-751", NewSike751(rand.Reader),
		"E1A758EC0D418BFE86D8077B5BB169133C06C1F2A067D8B202D9D058FFC51F63" +
			"FD26155A6577C74BA7F1A27E7BA51982517B923615DEB00BE408920A07831DF5" +
			"978CFDDD0BF690A264353A4A16B666F90586D7F89A193CE09375D389C1379A7A" +
			"528581C3ACB002CD2DC4F0FD672568FF9050BA8365C7FEFC5E6ED089B921DE68" +
			"04091A0744DE3EB14D426A3F7DA215C50312617C1C2697243980D06056F2CCE8" +
			"8AE7AE73C7343C0B7104C9F2870A94FED744CF6E94630514B6CEAB0E64733BB6" +
			"FA67B931E5D8206010475CBE8BC587248D65D89D8CD9C8BBFA93E8B5F9EB9130" +
			"773DED665D52ABBD91C4C8C255F73C0FC82501AE33330E9F308DE7177CBF83E4" +
			"E26E334D7CB09019E638147FC58ED372AF660F14C194BC80E9666325C98E0F80" +
			"877271D4A6BF514F603703D8A697874CD50A34D92F5AAEA84633CCF96801BD51" +
			"7BF425DEE4A32AAF06684052473EA14643C3D535440FB2240A988D09F297C5A3" +
			"88CB3DE60ED943F124034B90EFF611221F80F78EC124956338A105F6636B063D" +
			"7E48BFBD5D614310FB97D86F122E4AE6F9DDF4977A93ED7D0CE2A94E346A1A03" +
			"D3219CF21907B85A5BCDC713F93A4406A22E03B1655A66E1F6741A2F953E6FE0" +
			"868B2614BABEF1943BBBCB1B66D3E7017E533EA84F291240B56AB33EF1DC3F3D" +
			"E99DBF9E8BE51A0076E462BCDD825EA96D7F63C99177C305C257B31461F4C23D" +
			"43115F0220409E8880BBB2468586D03461E807BE824B693874911B2B52AF06FD" +
			"BDC47F5A0159729641A7C950AB9E03F2DC045135",
		"0001020304050607080900010203040506070809000102030405060708090102" +
			"8626ED79D451140800E03B59B956F8210E556067407D13DC90FA9E8B872BFB8F" +
			"AB0A7289852106E40538D3575C500201", MultiKEM{}},
}

// Encrypt, Decrypt, check if input/output plaintext is the same.
func testPKERoundTrip(t *testing.T, v sikeVec) {
	// Message to be encrypted
	var pt [common.MaxMsgBsz]byte
	var params = common.Params(v.id)
	var ct = make([]byte, v.kem.CiphertextSize())
	var msg = make([]byte, params.MsgLen)
	for i := range msg {
		msg[i] = byte(i)
	}

	// Import keys
	pkB := sidh.NewPublicKey(params.ID, sidh.KeyVariantSike)
	skB := sidh.NewPrivateKey(params.ID, sidh.KeyVariantSike)
	pkHex, err := hex.DecodeString(v.PkB)
	IsOk(t, err, "Test vector wrong")
	skHex, err := hex.DecodeString(v.PrB)
	IsOk(t, err, "Test vector wrong")
	err = pkB.Import(pkHex)
	IsOk(t, err, "Public key import failed")
	err = skB.Import(skHex)
	IsOk(t, err, "Private key import failed")
	err = v.kem.encrypt(ct, rng, pkB, msg[:])
	IsOk(t, err, "PKE roundtrip - encryption failed")
	ptLen, err := v.kem.decrypt(pt[:], skB, ct)
	IsOk(t, err, "PKE roundtrip - decription failed")

	if !bytes.Equal(pt[:ptLen], msg[:]) {
		t.Errorf("Decryption failed \n got : %X\n exp : %X", pt[:ptLen], msg)
	}
}

// Generate key and check if can encrypt.
func testPKEKeyGeneration(t *testing.T, v sikeVec) {
	var err error
	var params = common.Params(v.id)
	var pt [common.MaxMsgBsz]byte
	var msg = make([]byte, params.MsgLen)
	var ct = make([]byte, v.kem.CiphertextSize())
	// static buffer to ensure no overrides
	var pk = sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	var sk = sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)

	for i := range msg {
		msg[i] = byte(i)
	}

	err = sk.Generate(rng)
	IsOk(t, err, "PKE key generation")
	sk.GeneratePublicKey(pk)

	err = v.kem.encrypt(ct, rng, pk, msg[:])
	IsOk(t, err, "PKE encryption")
	ptLen, err := v.kem.decrypt(pt[:], sk, ct)
	IsOk(t, err, "PKE key decryption")

	if !bytes.Equal(pt[:ptLen], msg[:]) {
		t.Fatalf("Decryption failed \n got : %X\n exp : %X", pt, msg)
	}
}

func testNegativePKE(t *testing.T, v sikeVec) {
	var err error
	var msg [common.MaxMsgBsz]byte
	var ct = make([]byte, v.kem.CiphertextSize())
	var pk = sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	var sk = sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)

	// Generate key
	err = sk.Generate(rng)
	IsOk(t, err, "key generation")
	sk.GeneratePublicKey(pk)

	// bytelen(msg) - 1
	err = v.kem.encrypt(ct, rng, pk, msg[:v.kem.params.KemSize+8-1])
	if err == nil {
		t.Error(msg)
	}

	for _, v := range ct {
		if v != 0 {
			t.Fatal("Returned ciphertext must be not changed")
		}
	}
}

func testKEMRoundTrip(t *testing.T, pkB, skB []byte, v sikeVec) {
	// Import keys
	var err error
	var ssE [common.MaxSharedSecretBsz]byte
	var ssD [common.MaxSharedSecretBsz]byte
	var pk = sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	var sk = sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)
	var ct = make([]byte, v.kem.CiphertextSize())
	var ssBsz = v.kem.SharedSecretSize()

	err = pk.Import(pkB)
	IsOk(t, err, "Public key import failed")
	err = sk.Import(skB)
	IsOk(t, err, "Private key import failed")

	v.kem.Reset()
	err = v.kem.Encapsulate(ct, ssE[:], pk)
	IsOk(t, err, "Encapsulation failed")
	v.kem.Reset()
	err = v.kem.Decapsulate(ssD[:ssBsz], sk, pk, ct)
	IsOk(t, err, "Decapsulation failed")

	if !bytes.Equal(ssE[:v.kem.SharedSecretSize()], ssD[:v.kem.SharedSecretSize()]) {
		t.Errorf("Shared secrets from decapsulation and encapsulation differ [%s]", v.name)
	}
}

func testKEMKeyGeneration(t *testing.T, v sikeVec) {
	var ssE [common.MaxSharedSecretBsz]byte
	var ssD [common.MaxSharedSecretBsz]byte
	var ct = make([]byte, v.kem.CiphertextSize())

	sk := sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)
	pk := sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	IsOk(t, sk.Generate(rng), "error: key generation")
	sk.GeneratePublicKey(pk)

	// calculated shared secret
	v.kem.Reset()
	err := v.kem.Encapsulate(ct, ssE[:], pk)
	IsOk(t, err, "encapsulation failed")
	v.kem.Reset()
	err = v.kem.Decapsulate(ssD[:v.kem.SharedSecretSize()], sk, pk, ct)
	IsOk(t, err, "decapsulation failed")

	if !bytes.Equal(ssE[:], ssD[:]) {
		t.Fatalf("KEM failed \n encapsulated: %X\n decapsulated: %X", ssD[:], ssE[:])
	}
}

func testNegativeKEM(t *testing.T, v sikeVec) {
	var ssE [common.MaxSharedSecretBsz]byte
	var ssD [common.MaxSharedSecretBsz]byte
	var ssTmp [common.MaxSharedSecretBsz]byte
	var ct = make([]byte, v.kem.CiphertextSize())
	var ssBsz = v.kem.SharedSecretSize()

	sk := sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)
	pk := sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	IsOk(t, sk.Generate(rng), "error: key generation")
	sk.GeneratePublicKey(pk)

	v.kem.Reset()
	err := v.kem.Encapsulate(ct, ssE[:], pk)
	IsOk(t, err, "pre-requisite for a test failed")

	// Try decapsulate too small ciphertext
	v.kem.Reset()
	IsOk(
		t,
		testPanic(func() { _ = v.kem.Decapsulate(ssTmp[:ssBsz], sk, pk, ct[:len(ct)-2]) }),
		"Decapsulation must panic if ciphertext is too small")

	ctTmp := make([]byte, len(ct)+1)
	// Try decapsulate too big ciphertext
	v.kem.Reset()
	IsOk(
		t,
		testPanic(func() { _ = v.kem.Decapsulate(ssTmp[:ssBsz], sk, pk, ctTmp) }),
		"Decapsulation must panic if ciphertext is too big")

	// Change ciphertext
	ct[0] = ct[0] - 1
	v.kem.Reset()
	err = v.kem.Decapsulate(ssD[:ssBsz], sk, pk, ct)
	IsOk(t, err, "decapsulation returns error when invalid ciphertext provided")

	if bytes.Equal(ssE[:], ssD[:]) {
		// no idea how this could ever happen, but it would be very bad
		t.Error("critical error")
	}

	// Try encapsulating with SIDH key
	pkSidh := sidh.NewPublicKey(v.id, sidh.KeyVariantSidhB)
	prSidh := sidh.NewPrivateKey(v.id, sidh.KeyVariantSidhB)
	v.kem.Reset()
	IsOk(
		t,
		testPanic(func() { _ = v.kem.Encapsulate(ct, ssE[:], pkSidh) }),
		"encapsulation accepts SIDH public key")

	// Try decapsulating with SIDH key
	v.kem.Reset()
	IsOk(
		t,
		testPanic(func() { _ = v.kem.Decapsulate(ssD[:ssBsz], prSidh, pk, ct) }),
		"encapsulation accepts SIDH public key")
}

// In case invalid ciphertext is provided, SIKE's decapsulation must
// return same (but unpredictable) result for a given key.
func testNegativeKEMSameWrongResult(t *testing.T, v sikeVec) {
	var ssE [common.MaxSharedSecretBsz]byte
	var ssD1 [common.MaxSharedSecretBsz]byte
	var ssD2 [common.MaxSharedSecretBsz]byte
	var ct = make([]byte, v.kem.CiphertextSize())
	var ssBsz = v.kem.SharedSecretSize()

	sk := sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)
	pk := sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	IsOk(t, sk.Generate(rng), "error: key generation")
	sk.GeneratePublicKey(pk)

	v.kem.Reset()
	err := v.kem.Encapsulate(ct, ssE[:], pk)
	IsOk(t, err, "pre-requisite for a test failed")

	// make ciphertext wrong
	ct[0] = ct[0] - 1
	v.kem.Reset()
	err = v.kem.Decapsulate(ssD1[:ssBsz], sk, pk, ct)
	IsOk(t, err, "pre-requisite for a test failed")

	// change secret keysecond decapsulation must be done with same, but imported private key
	var expSk [common.MaxSikePrivateKeyBsz]byte
	sk.Export(expSk[:])

	// create new private key
	sk = sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)
	IsOk(t, sk.Import(expSk[:sk.Size()]), "import failed")

	// try decapsulating again.
	v.kem.Reset()
	err = v.kem.Decapsulate(ssD2[:ssBsz], sk, pk, ct)
	IsOk(t, err, "pre-requisite for a test failed")

	// ssD1 must be same as ssD2
	if !bytes.Equal(ssD1[:], ssD2[:]) {
		t.Error("decapsulation is insecure")
	}

	// ssD1 and ssD2 must be different than ssE
	if bytes.Equal(ssE[:], ssD1[:]) || bytes.Equal(ssE[:], ssD2[:]) {
		// this test requires that decapsulation returns wrong result
		t.Errorf("test implementation error")
	}
}

// Test MultiKEM
func TestMultiKemRoundTrip(t *testing.T) {
	var err error
	var mkem MultiKEM
	var pks []*sidh.PublicKey
	var sks []*sidh.PrivateKey
	var ss_out [common.MaxSharedSecretBsz]byte
	var ss_in [common.MaxSharedSecretBsz]byte

	mkem.Allocate(common.Fp434, 10, rng)
	pks = make([]*sidh.PublicKey, len(mkem.cts))
	sks = make([]*sidh.PrivateKey, len(mkem.cts))
	// create recipients keys

	for i, _ := range mkem.cts {
		pks[i] = mkem.NewPublicKey()
		sks[i] = mkem.NewPrivateKey()
		err = sks[i].Generate(rng)
		IsOk(t, err, "PKE key generation")
		sks[i].GeneratePublicKey(pks[i])
	}

	err = mkem.Encapsulate(ss_out[:], pks)
	IsOk(t, err, "Multi KEM failed")
	for i := 0; i < len(mkem.cts); i++ {
		err = mkem.Decapsulate(ss_in[:], sks[i], pks[i], mkem.cts[i][:mkem.KemSize()])
		IsOk(t, err, "Decaps failed")
		Ok(t, bytes.Equal(ss_out[:mkem.KemSize()], ss_in[:mkem.KemSize()]), "shared secret equal")
	}
}

// Interface to "testing"

/* -------------------------------------------------------------------------
   Wrappers for 'testing' SIDH
   -------------------------------------------------------------------------*/
func testSike(t *testing.T, m *map[uint8]sikeVec, f func(*testing.T, sikeVec)) {
	for i := range *m {
		v := (*m)[i]
		t.Run(v.name, func(t *testing.T) { f(t, v) })
	}
}

func TestPKERoundTrip(t *testing.T)     { testSike(t, &tdataSike, testPKERoundTrip) }
func TestPKEKeyGeneration(t *testing.T) { testSike(t, &tdataSike, testPKEKeyGeneration) }
func TestNegativePKE(t *testing.T)      { testSike(t, &tdataSike, testNegativePKE) }
func TestKEMKeyGeneration(t *testing.T) { testSike(t, &tdataSike, testKEMKeyGeneration) }
func TestNegativeKEM(t *testing.T)      { testSike(t, &tdataSike, testNegativeKEM) }
func TestNegativeKEMSameWrongResult(t *testing.T) {
	testSike(t, &tdataSike, testNegativeKEMSameWrongResult)
}

func TestKEMRoundTrip(t *testing.T) {
	for _, val := range tdataSike {
		//		fmt.Printf("\tTesting: %s\n", val.name)
		pk, err := hex.DecodeString(val.PkB)
		IsOk(t, err, "public key B not a number")
		sk, err := hex.DecodeString(val.PrB)
		IsOk(t, err, "private key B not a number")
		testKEMRoundTrip(t, pk, sk, val)
	}
}

/* -------------------------------------------------------------------------
   Benchmarking
   -------------------------------------------------------------------------*/

func benchSike(t *testing.B, m *map[uint8]sikeVec, f func(*testing.B, sikeVec)) {
	for i := range *m {
		v := (*m)[i]
		t.Run(v.name, func(t *testing.B) { f(t, v) })
	}
}

func benchKeygen(b *testing.B, v sikeVec) {
	pub := sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	prv := sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)
	_ = prv.Generate(rng)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		prv.GeneratePublicKey(pub)
	}
}

func benchmarkEncaps(b *testing.B, v sikeVec) {
	pub := sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	prv := sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)

	if prv.Generate(rng) != nil {
		b.FailNow()
	}
	prv.GeneratePublicKey(pub)

	var ct [common.MaxCiphertextBsz]byte
	var ss [common.MaxSharedSecretBsz]byte

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.kem.Reset()
		_ = v.kem.Encapsulate(ct[:], ss[:], pub)
	}
}

func benchmarkMultiEncaps_100keys(b *testing.B, v sikeVec) {
	var ss [common.MaxSharedSecretBsz]byte
	v.mkem.Allocate(v.id, 100, rng)
	// list of public keys
	pks := make([]*sidh.PublicKey, len(v.mkem.cts))
	sk := v.mkem.NewPrivateKey()

	// create keys
	for i, _ := range v.mkem.cts {
		pks[i] = v.mkem.NewPublicKey()
		_ = sk.Generate(rng)
		sk.GeneratePublicKey(pks[i])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.mkem.Encapsulate(ss[:], pks)
	}
}

func benchmarkDecaps(b *testing.B, v sikeVec) {
	var ct [common.MaxCiphertextBsz]byte
	var ss [common.MaxSharedSecretBsz]byte
	var ssBsz = v.kem.SharedSecretSize()

	pkA := sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	prvA := sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)
	pkB := sidh.NewPublicKey(v.id, sidh.KeyVariantSike)
	prvB := sidh.NewPrivateKey(v.id, sidh.KeyVariantSike)

	if prvA.Generate(rng) != nil || prvB.Generate(rng) != nil {
		b.FailNow()
	}

	prvA.GeneratePublicKey(pkA)
	prvB.GeneratePublicKey(pkB)

	v.kem.Reset()
	err := v.kem.Encapsulate(ct[:], ss[:], pkA)
	if err != nil {
		b.FailNow()
	}

	ctSlc := ct[:v.kem.CiphertextSize()]

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.kem.Reset()
		_ = v.kem.Decapsulate(ss[:ssBsz], prvA, pkB, ctSlc)
	}
}

func BenchmarkKeygen(b *testing.B) { benchSike(b, &tdataSike, benchKeygen) }
func BenchmarkEncaps(b *testing.B) { benchSike(b, &tdataSike, benchmarkEncaps) }
func BenchmarkMultiEncaps_100keys(b *testing.B) {
	benchSike(b, &tdataSike, benchmarkMultiEncaps_100keys)
}
func BenchmarkDecaps(b *testing.B) { benchSike(b, &tdataSike, benchmarkDecaps) }

func ExampleKEM() {
	// Allice's key pair
	prvA := sidh.NewPrivateKey(sidh.Fp503, sidh.KeyVariantSike)
	pubA := sidh.NewPublicKey(sidh.Fp503, sidh.KeyVariantSike)
	// Bob's key pair
	prvB := sidh.NewPrivateKey(sidh.Fp503, sidh.KeyVariantSike)
	pubB := sidh.NewPublicKey(sidh.Fp503, sidh.KeyVariantSike)
	// Generate keypair for Allice
	err := prvA.Generate(rng)
	if err != nil {
		panic(err)
	}
	prvA.GeneratePublicKey(pubA)
	// Generate keypair for Bob
	err = prvB.Generate(rng)
	if err != nil {
		panic(err)
	}
	prvB.GeneratePublicKey(pubB)
	// Initialize internal KEM structures
	var kem = NewSike503(rng)
	// Create buffers for ciphertext, shared secret received
	// from encapsulation and shared secret from decapsulation
	ct := make([]byte, kem.CiphertextSize())
	ssE := make([]byte, kem.SharedSecretSize())
	ssD := make([]byte, kem.SharedSecretSize())
	// Allice performs encapsulation with Bob's public key
	err = kem.Encapsulate(ct, ssE, pubB)
	if err != nil {
		panic(err)
	}
	// Bob performs decapsulation with his key pair
	err = kem.Decapsulate(ssD, prvB, pubB, ct)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", bytes.Equal(ssE, ssD))

	// Bob performs encapsulation with Allices's public key
	err = kem.Encapsulate(ct, ssE, pubA)
	if err != nil {
		panic(err)
	}
	// Allice performs decapsulation with hers key pair
	err = kem.Decapsulate(ssD, prvA, pubA, ct)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", bytes.Equal(ssE, ssD))

	// Output:
	// true
	// true
}
