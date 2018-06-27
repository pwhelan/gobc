package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"
)

var chain = `
{
  "Hash": "75d91be3fe7f83e50e7036aba4383b601404734db31ffffbddc215873289c78427fff7ca67fd4c83e79793b8cbb6440ee46760affadf48fc16473ee3e6424bb5",
  "Index": 0,
  "Timestamp": 1527949914,
  "Difficulty": 0,
  "CumulativeDifficulty": 0,
  "PrevHash": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "Transactions": null,
  "Nonce": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
}
{
  "Hash": "2913a65073a3b42ac489c980430973fc63fd6f075a71c27e8a142adaf4f7e774613da6cbe482273acf37b5352c0e752070f007dfdc09ba985801f8cc25152d75",
  "Index": 1,
  "Timestamp": 1529426969,
  "Difficulty": 1,
  "CumulativeDifficulty": 1,
  "PrevHash": "75d91be3fe7f83e50e7036aba4383b601404734db31ffffbddc215873289c78427fff7ca67fd4c83e79793b8cbb6440ee46760affadf48fc16473ee3e6424bb5",
  "Transactions": null,
  "Nonce": "ca6d6ec73b725a5b96a0bf3f4f76d150e8ba41d9d9130bee4a6055effc72d968d560f08bc41adca7471e81233276c2388beedc74d219caceb699b5d1b811cd78"
}
{
  "Hash": "d3edda3677c26c405938839059899b62603812f818e03c4ab701f8695432ade37cdbcc9aaff0560f2554925fa62028a115fa12db01f3e9433fbf75fe20de247e",
  "Index": 2,
  "Timestamp": 1529426969,
  "Difficulty": 0,
  "CumulativeDifficulty": 1,
  "PrevHash": "2913a65073a3b42ac489c980430973fc63fd6f075a71c27e8a142adaf4f7e774613da6cbe482273acf37b5352c0e752070f007dfdc09ba985801f8cc25152d75",
  "Transactions": null,
  "Nonce": "5c9ba0a9bce479a2f7e87e925c8a8c6048141aa46906f013134da34b995b30fc0fafd4f4a8ad2247d95ca7290b10cc7b695bcabfed0ab31d195200a63e427cd7"
}
{
  "Hash": "1e247349ea8de63369dcb8b32bda0e1c6c0ed8faebc412486a7819f8ca8355a31d896498ce157dc4a7448c7d2a7969a1c0d99bdf908739548b248fa38372a4ad",
  "Index": 3,
  "Timestamp": 1529426969,
  "Difficulty": 2,
  "CumulativeDifficulty": 3,
  "PrevHash": "d3edda3677c26c405938839059899b62603812f818e03c4ab701f8695432ade37cdbcc9aaff0560f2554925fa62028a115fa12db01f3e9433fbf75fe20de247e",
  "Transactions": null,
  "Nonce": "2a9773296be364aaccd8cf3b257dd931c93f3810116c4a80ec021da6518f44d75bbd202b97ad64d9dca90db80de121db11ea6e6465cc20be343a24fa8a43e930"
}
{
  "Hash": "210c1bed1627cd8a93edeea21c1c5d20ec47152cf9bd49318cc43ac97feed68292f81babafef0b9704b29e41ae2dec72d9ff1c19a0d754568ec964637a1bb745",
  "Index": 4,
  "Timestamp": 1529426969,
  "Difficulty": 1,
  "CumulativeDifficulty": 4,
  "PrevHash": "1e247349ea8de63369dcb8b32bda0e1c6c0ed8faebc412486a7819f8ca8355a31d896498ce157dc4a7448c7d2a7969a1c0d99bdf908739548b248fa38372a4ad",
  "Transactions": null,
  "Nonce": "bc39cd2e75ea4bc89682be04b6561e8b599ca4337295cf378bbe9ee3db66bedde974d145e8704c283e21c1440c391715f7cfc48ab86b67b5ffbc716ca2d671a9"
}
{
  "Hash": "140acd9a7098a8e7387aa2812ebb725ba5c9dbdfec3684845b72d7ac26e256e94ebaac542bf485490eb62ded823c2ec782b8744a1afe04fec1e033d560dfc5ab",
  "Index": 5,
  "Timestamp": 1529426969,
  "Difficulty": 3,
  "CumulativeDifficulty": 7,
  "PrevHash": "210c1bed1627cd8a93edeea21c1c5d20ec47152cf9bd49318cc43ac97feed68292f81babafef0b9704b29e41ae2dec72d9ff1c19a0d754568ec964637a1bb745",
  "Transactions": null,
  "Nonce": "ef0627b3574d9aa489061069489643ae070b75d716d4e68bfb4b17e1fdcbe6ddde4c88487e11a8fad4e1954d70204e44e5e0006c68da3b48fe0607ceab7cc18c"
}
{
  "Hash": "17e28266d32951372567cb939e4c53c02d72076d5e4b779edca5344c3e19827d0b33abeec4d9117ac10f1583938ebae805e7cd2eff0e920168d85bfacc4f3dfd",
  "Index": 6,
  "Timestamp": 1529426969,
  "Difficulty": 2,
  "CumulativeDifficulty": 9,
  "PrevHash": "140acd9a7098a8e7387aa2812ebb725ba5c9dbdfec3684845b72d7ac26e256e94ebaac542bf485490eb62ded823c2ec782b8744a1afe04fec1e033d560dfc5ab",
  "Transactions": null,
  "Nonce": "cea7b265d95fadd6df9bac35fa0ed94252b387dc4323fc9abbb387c44ab51ba98bc3dbc7ed946b2cdc5dd8f2a42d3e0cf9c7ebc348a92212f648c4c4cd77bc98"
}
{
  "Hash": "0879f5e88eb86c132602aef9781997821eff2c4d8b89fa6abaaf0d9fe491a6157089a1ff8814a6e49acf8c518526ea512dfd6b8740b452376e4298aef2ae0901",
  "Index": 7,
  "Timestamp": 1529426969,
  "Difficulty": 4,
  "CumulativeDifficulty": 13,
  "PrevHash": "17e28266d32951372567cb939e4c53c02d72076d5e4b779edca5344c3e19827d0b33abeec4d9117ac10f1583938ebae805e7cd2eff0e920168d85bfacc4f3dfd",
  "Transactions": null,
  "Nonce": "911a2b7b23cca6c7ca4d7c654fc8f597cf9bd8692e668e8b9428f96763fa6458addb975129d094031c8c34da0b7317024e4b982276e4b47eacd47f0e38c4716c"
}
{
  "Hash": "0932dbb37ae4a06ab48cb60212a845c824c0721257d7529b69ce9e588fab11ab52926116bf03b2bf2b4a5edbc7580b563cbf4907da0e3860ddbd2284cb8cc666",
  "Index": 8,
  "Timestamp": 1529426969,
  "Difficulty": 3,
  "CumulativeDifficulty": 16,
  "PrevHash": "0879f5e88eb86c132602aef9781997821eff2c4d8b89fa6abaaf0d9fe491a6157089a1ff8814a6e49acf8c518526ea512dfd6b8740b452376e4298aef2ae0901",
  "Transactions": null,
  "Nonce": "351123fc675e93be1f3147e9a71bcb322fa2e082833a4032e8ca089a7f886f2e233da72f6866cb593c196ac710047f8b3e709f62c189fe30f94f0951c36344cd"
}
{
  "Hash": "074a2b7b668330ba0133e7d6bb9eadd29d4529ffa58c2a114a3be61600c0a7c8d6a675d5bee6abfe776cf664a12d480ba38c18e513f217f57eda8ac2d996b4f8",
  "Index": 9,
  "Timestamp": 1529426969,
  "Difficulty": 5,
  "CumulativeDifficulty": 21,
  "PrevHash": "0932dbb37ae4a06ab48cb60212a845c824c0721257d7529b69ce9e588fab11ab52926116bf03b2bf2b4a5edbc7580b563cbf4907da0e3860ddbd2284cb8cc666",
  "Transactions": null,
  "Nonce": "3dd8146899c21bbef56d2d049a4608d97614e58e1323d0707f4e1c5910863e6243fe43385aecdc616b009a6514eb44e8233f84013e3b8eb5369814d0f2d5519b"
}
{
  "Hash": "008aff24340bbe85d0dc4debf83f85f4bf9107a085a2055da7e5a6bf681e42fcf7e91269ea71a5bbba6e7e54870a8c2040cc91bd4dc6f0f1447601dc3c81140c",
  "Index": 10,
  "Timestamp": 1529426969,
  "Difficulty": 4,
  "CumulativeDifficulty": 25,
  "PrevHash": "074a2b7b668330ba0133e7d6bb9eadd29d4529ffa58c2a114a3be61600c0a7c8d6a675d5bee6abfe776cf664a12d480ba38c18e513f217f57eda8ac2d996b4f8",
  "Transactions": null,
  "Nonce": "eb95804902b31e8a79fd2cf6f47522a86e6f3b0c215eab4e740c4744ae19b6a1e96cec2e6fb66c3b6598611bfa3ddf1a468c3fa8626d7b0d9e7a83b2a6d76b66"
}
{
  "Hash": "03f93f3710678527b86a70abdf32be1ecfc46b38b353bfb4445c85d501df09c9ed0ab320291cb8f3d78e9429ee354c71ecbfb2b232936c4219ff717ee5cfdd0c",
  "Index": 11,
  "Timestamp": 1529426969,
  "Difficulty": 6,
  "CumulativeDifficulty": 31,
  "PrevHash": "008aff24340bbe85d0dc4debf83f85f4bf9107a085a2055da7e5a6bf681e42fcf7e91269ea71a5bbba6e7e54870a8c2040cc91bd4dc6f0f1447601dc3c81140c",
  "Transactions": null,
  "Nonce": "58f916013db9d9e4feaddc1b9f6c8aebe27c667af03f4086fde13aee10582f1c4f46c118c2952654e21d7f695454fe5559b5abda3087a524c6b9488db615a856"
}
{
  "Hash": "07528af953d2e3ff6a38075c4e6ab716258e43b6739200ae47fdfaf7e99e1900d7541489255efb06f0ceaa0603725857ae32faace65c6e731fe8e57bc40c238b",
  "Index": 12,
  "Timestamp": 1529426969,
  "Difficulty": 5,
  "CumulativeDifficulty": 36,
  "PrevHash": "03f93f3710678527b86a70abdf32be1ecfc46b38b353bfb4445c85d501df09c9ed0ab320291cb8f3d78e9429ee354c71ecbfb2b232936c4219ff717ee5cfdd0c",
  "Transactions": null,
  "Nonce": "dc3e70ede2353a17762963eebef4a2d4478908a8feb725a36706ce3f9808732b1a5241cbe92f6f5d664cacbaa3936c174f8a8c684515c0dc8271f1f415cc6aed"
}
{
  "Hash": "00f1f5b5dabb3bade7539d238d36fa04d54313d75e150a4d11e55ee44cdf73d80aaa54dd7b8802d6129e217ea10932ca63611d84b91f36239674c2235d57cfac",
  "Index": 13,
  "Timestamp": 1529426969,
  "Difficulty": 7,
  "CumulativeDifficulty": 43,
  "PrevHash": "07528af953d2e3ff6a38075c4e6ab716258e43b6739200ae47fdfaf7e99e1900d7541489255efb06f0ceaa0603725857ae32faace65c6e731fe8e57bc40c238b",
  "Transactions": null,
  "Nonce": "c825088f397411e9ad18bf26e234faba188bfb8faaac46088dbe4b1a052b9be22710feeb58fcdd32f5770cd08b8b71f39b5d38cd0f6a435d22e8d6be570b031a"
}
{
  "Hash": "03f2996fb1aee2c67f8190fa3d46654393bfac1d8a434123b56341f7a0cac32a4bb2ae32a76415115e5552c02222d608d36e451025b38a51fcdba48478b53b1d",
  "Index": 14,
  "Timestamp": 1529426969,
  "Difficulty": 6,
  "CumulativeDifficulty": 49,
  "PrevHash": "00f1f5b5dabb3bade7539d238d36fa04d54313d75e150a4d11e55ee44cdf73d80aaa54dd7b8802d6129e217ea10932ca63611d84b91f36239674c2235d57cfac",
  "Transactions": null,
  "Nonce": "d444a7e2e40feea314c112e924d2f63352a699236a95653f0002f959c3a86b11a5bc3810583d368d2b5c486201894d3563ebd4bdc7cbddde44130a136f52f9c4"
}
{
  "Hash": "00d3ebdbd6d05f9ce6df0277acb07bd16f16cd7786894ced3180101d927c6993495a405e455ab55b4f93141555304051bde663b48d45c7ee496fd79f0b4c5cdb",
  "Index": 15,
  "Timestamp": 1529426969,
  "Difficulty": 8,
  "CumulativeDifficulty": 57,
  "PrevHash": "03f2996fb1aee2c67f8190fa3d46654393bfac1d8a434123b56341f7a0cac32a4bb2ae32a76415115e5552c02222d608d36e451025b38a51fcdba48478b53b1d",
  "Transactions": null,
  "Nonce": "ea15ae62ffab9cce5ec2af932f050a97c38dcfdf9cb4ee6f4d219d71c20acc28203271aa6c2afca8efbb68fd0852581aa69971c31c2eb93e853e821485dcb7fe"
}
{
  "Hash": "004dc5902dce5b8bec6104d9eef4550a71bd07d9d983d2e808f301fe99794044a4872a7139946b855fdc0b3a0779617c079f58eb738d6bddeb46e63444c1432d",
  "Index": 16,
  "Timestamp": 1529426969,
  "Difficulty": 7,
  "CumulativeDifficulty": 64,
  "PrevHash": "00d3ebdbd6d05f9ce6df0277acb07bd16f16cd7786894ced3180101d927c6993495a405e455ab55b4f93141555304051bde663b48d45c7ee496fd79f0b4c5cdb",
  "Transactions": null,
  "Nonce": "d5047882bb2308bbf79506a2f1802841b3751115e7b31dc59d84df9a84fcc7d8bd1a98c3618f8b83f7fbaf26012a0d0db6b46a6245c759d6cd907a9a86244b38"
}
{
  "Hash": "000c79494c2d3ddf7f70e1cf7ff0c2d463bec1c44e0c7820b1e0e1ca7d20d00394a1243608600a335374031dc28246534db4ea42fc5933b2692f0c36d82d88b9",
  "Index": 17,
  "Timestamp": 1529426969,
  "Difficulty": 9,
  "CumulativeDifficulty": 73,
  "PrevHash": "004dc5902dce5b8bec6104d9eef4550a71bd07d9d983d2e808f301fe99794044a4872a7139946b855fdc0b3a0779617c079f58eb738d6bddeb46e63444c1432d",
  "Transactions": null,
  "Nonce": "c1f30895263f781158c50ef12a87f99d04507067d7741fb26b1de027334ea7b84144a77ba23e6cf5d31c38454f1ffcbc10247a04fd24dd096038eed716e2f7cf"
}
{
  "Hash": "00ec14e2ca98f23dafd389ed3351837816cb9034056b92dc231e0cbaee1157658fed550b936e3cb6c2cbfd8c6f2868ae570832a38a2979e3a603f346e0defdcd",
  "Index": 18,
  "Timestamp": 1529426969,
  "Difficulty": 8,
  "CumulativeDifficulty": 81,
  "PrevHash": "000c79494c2d3ddf7f70e1cf7ff0c2d463bec1c44e0c7820b1e0e1ca7d20d00394a1243608600a335374031dc28246534db4ea42fc5933b2692f0c36d82d88b9",
  "Transactions": null,
  "Nonce": "c35226c63063553396c991b4f4cceb9111dd1b9443caeaef777245008d5f4da30d06475f30cf2c2d6685f7e13d0083d20150cb6b5529896d3b3e3258d16a45c3"
}
{
  "Hash": "0010c025d7287f2844d03c9d9ea71b5289e7ff791aea741fdb548f2f0db319dc042bb66da53003c910a490514572cbf9660ba6a5f8b316e896e77b6abd7e8c43",
  "Index": 19,
  "Timestamp": 1529426969,
  "Difficulty": 10,
  "CumulativeDifficulty": 91,
  "PrevHash": "00ec14e2ca98f23dafd389ed3351837816cb9034056b92dc231e0cbaee1157658fed550b936e3cb6c2cbfd8c6f2868ae570832a38a2979e3a603f346e0defdcd",
  "Transactions": null,
  "Nonce": "9b9a6837e86ccf10dd3b2c2b67dfabf2cc7503e26b63f28d7a081133edc85ab4d564023139334a2c64c5bfbca4b1fb2d6f91ad99c140b363348584d7ebc731f4"
}
{
  "Hash": "007d59e70509d7b1b8e16dd50448cd5f2550009b109a4cf0e1d0aef0cbe40a1e53e20a6b63ac79d7c403b96c21b5357a0af4d926fc4dc5878198332155cc5afd",
  "Index": 20,
  "Timestamp": 1529426969,
  "Difficulty": 9,
  "CumulativeDifficulty": 100,
  "PrevHash": "0010c025d7287f2844d03c9d9ea71b5289e7ff791aea741fdb548f2f0db319dc042bb66da53003c910a490514572cbf9660ba6a5f8b316e896e77b6abd7e8c43",
  "Transactions": null,
  "Nonce": "407cbe243a4d208003349785a0a4ef7ab4034cdfb81930eb594f3e72d781dfeb4bcea5156b4778ee6d319437885af30afd6ae974742a8d644876cddf8f9f3b90"
}
{
  "Hash": "000983b6cb74607cd0e060d7f602e1c79d242457097815717e0def7c7aa45c93648aa30391e74bab5be039140db267ec3642ded8325e241a603a0b17809f7708",
  "Index": 21,
  "Timestamp": 1529426969,
  "Difficulty": 11,
  "CumulativeDifficulty": 111,
  "PrevHash": "007d59e70509d7b1b8e16dd50448cd5f2550009b109a4cf0e1d0aef0cbe40a1e53e20a6b63ac79d7c403b96c21b5357a0af4d926fc4dc5878198332155cc5afd",
  "Transactions": null,
  "Nonce": "299eebfb9f6ced150831a74fc84659184c22680ad8ac3f7bdf21f498b65ef57409987e3de990989f6daad881a9c8e782b72c97ec3e1114cb5c7a426a176e0faf"
}
{
  "Hash": "001168a49baa94bcbb9094cc74d647dc7b7a4eff840066069c4e9319bcc8ff888f25f5d494bd8a94f289b4cdf65f8df0b11ecfcbb364fc9135dcc95e1b19db70",
  "Index": 22,
  "Timestamp": 1529426969,
  "Difficulty": 10,
  "CumulativeDifficulty": 121,
  "PrevHash": "000983b6cb74607cd0e060d7f602e1c79d242457097815717e0def7c7aa45c93648aa30391e74bab5be039140db267ec3642ded8325e241a603a0b17809f7708",
  "Transactions": null,
  "Nonce": "3fccdf1420b3e7a22a5ffd16756d6d975d203a67c4b36d86ef3440ce2d9cd40f156a758cd844fcd5b854097445d4b6db81c6609eb29b65a1a470e96d8b12243a"
}
{
  "Hash": "00061819e11738e04e8c11e965efbad01750bbb5c95270ce8f8b3ce1955cdd5218f2510685bdb0ff49358af58a446709c2b02f03d7f2acb08519543701a3562a",
  "Index": 23,
  "Timestamp": 1529426969,
  "Difficulty": 12,
  "CumulativeDifficulty": 133,
  "PrevHash": "001168a49baa94bcbb9094cc74d647dc7b7a4eff840066069c4e9319bcc8ff888f25f5d494bd8a94f289b4cdf65f8df0b11ecfcbb364fc9135dcc95e1b19db70",
  "Transactions": null,
  "Nonce": "b92d7f2d2d41973a233985224e757ccc333e1a08474e41d59705d36817343029876a013e32c8ae42ed1795832f815aaa0c74bd06e8fff71d43f5f8d91aafdd2a"
}
{
  "Hash": "0012c539b9dec47c320bf083dc4cca0517aa081db63d7b6a3626f86bd0f415d5c9950c28e3bbeb6d97a5a1145b37f15f8c6460d0417b4b9470c671b7c0c46e58",
  "Index": 24,
  "Timestamp": 1529426969,
  "Difficulty": 11,
  "CumulativeDifficulty": 144,
  "PrevHash": "00061819e11738e04e8c11e965efbad01750bbb5c95270ce8f8b3ce1955cdd5218f2510685bdb0ff49358af58a446709c2b02f03d7f2acb08519543701a3562a",
  "Transactions": null,
  "Nonce": "328bfc02b79d2623c0e25554939ef074bc5a849b94a6d3181d436890286e9f364fe657fd3755df9d2104c37536a0a3e97fd8cda0fba94881070545efb59c94e3"
}
{
  "Hash": "0005144dd7c9a1b8c7aa74a161973a20396008fb1116e86f1c9781e21c904865049cff4314ae0ac9eeb67d84f59d2e588b28a41bd295f9ce0464b86ed68394df",
  "Index": 25,
  "Timestamp": 1529426969,
  "Difficulty": 13,
  "CumulativeDifficulty": 157,
  "PrevHash": "0012c539b9dec47c320bf083dc4cca0517aa081db63d7b6a3626f86bd0f415d5c9950c28e3bbeb6d97a5a1145b37f15f8c6460d0417b4b9470c671b7c0c46e58",
  "Transactions": null,
  "Nonce": "1a8c6a1e89ca61372ecada14acc1eefe4f1a6ea4836725e5450d8e502f2b181ae4be488f11970bd8d31f9d89c7cedb9c3988d75694b6bbb684579f643af53def"
}
{
  "Hash": "0007ed8f4d1868c1c0e8c85d550dda2d139aa1e6f97b66fc747120987b9b04b7c58443605756ab342c68096b83b17c9cd5a3694b1f8d95e105d2faf01be18e7a",
  "Index": 26,
  "Timestamp": 1529426969,
  "Difficulty": 12,
  "CumulativeDifficulty": 169,
  "PrevHash": "0005144dd7c9a1b8c7aa74a161973a20396008fb1116e86f1c9781e21c904865049cff4314ae0ac9eeb67d84f59d2e588b28a41bd295f9ce0464b86ed68394df",
  "Transactions": null,
  "Nonce": "cdf2cbf1a734b5d6d26e68cd342ca5d00e6e38937704f9b2664653d4edd0744d9d85878dcad8c1a9d11ab95823b4f1fbcfa1c917dcbb30b873f976415c10f471"
}
{
  "Hash": "0000d3aadf873bab7558a8d35cbe84430df9ab451649dbe91cfb6b1fd2825c1e3580be94fe07a17ffaaa309503546cd0f98cdda628f9dbab936aa1bf141f45db",
  "Index": 27,
  "Timestamp": 1529426969,
  "Difficulty": 14,
  "CumulativeDifficulty": 183,
  "PrevHash": "0007ed8f4d1868c1c0e8c85d550dda2d139aa1e6f97b66fc747120987b9b04b7c58443605756ab342c68096b83b17c9cd5a3694b1f8d95e105d2faf01be18e7a",
  "Transactions": null,
  "Nonce": "cd3ecd81b792b91dce6df05ba535a435e1552006ece2405976031e833f67656ba5490ac13983e303f699b2ab15ccf7aa9aa192501b97086311833328b44a3006"
}
{
  "Hash": "0003538cc4b4cc6446a779840a75eaa3361996e76ba6e4b0c9c07f0a8203a4d4a3021c448b637c70231a3283d5602724e88decd6b37a36dea34e288faaca84ef",
  "Index": 28,
  "Timestamp": 1529426970,
  "Difficulty": 13,
  "CumulativeDifficulty": 196,
  "PrevHash": "0000d3aadf873bab7558a8d35cbe84430df9ab451649dbe91cfb6b1fd2825c1e3580be94fe07a17ffaaa309503546cd0f98cdda628f9dbab936aa1bf141f45db",
  "Transactions": null,
  "Nonce": "bd21bd314937fefc28aabd348e7856f7c3e4b84336ad2c5e3db7eed08ae4f7b7d5ed924253eab0920efadcd12930471fd316aba7c18ff8590054d7854943a6ae"
}
{
  "Hash": "0000384b9e1b89c966801869c6e68932d0f2dfca1c2cb7eadc3b8658aa95d73791cc5546119822ca8695f78d7c33587be820f2fe131738809b5844d7472ffc09",
  "Index": 29,
  "Timestamp": 1529426970,
  "Difficulty": 15,
  "CumulativeDifficulty": 211,
  "PrevHash": "0003538cc4b4cc6446a779840a75eaa3361996e76ba6e4b0c9c07f0a8203a4d4a3021c448b637c70231a3283d5602724e88decd6b37a36dea34e288faaca84ef",
  "Transactions": null,
  "Nonce": "c4ee64bf275f74fab7daba01f77a38c6177269d000e993752a1e90f3867348e39b88cb90db3d9b041271d17b53a61161c6c5161f615249cd5f2588ed8fc5251f"
}
{
  "Hash": "0002ece01dccee9030181cddce8f35d0f95f82d70ff59d6e5ae78325f5a5cc9b9186be67efed24048cb6e6ad38605217242023a2745209e1ed7479cba400e6a2",
  "Index": 30,
  "Timestamp": 1529426970,
  "Difficulty": 14,
  "CumulativeDifficulty": 225,
  "PrevHash": "0000384b9e1b89c966801869c6e68932d0f2dfca1c2cb7eadc3b8658aa95d73791cc5546119822ca8695f78d7c33587be820f2fe131738809b5844d7472ffc09",
  "Transactions": null,
  "Nonce": "6781a816ff0be64a55b983267f6eec2e23a66d4707a31eb183522e7d587bb861578b596a1660759d500e193c3fb1bd41e1617e2c2b6dee449fee457e46de6dee"
}
{
  "Hash": "00001d3c038a50185c39d46cbadc3da102779f4a8065a776cbf858b44483a5532095ff88c291bf154a02a6df0ac34c210ca931fa1c635c072ac7b66db06267bd",
  "Index": 31,
  "Timestamp": 1529426970,
  "Difficulty": 16,
  "CumulativeDifficulty": 241,
  "PrevHash": "0002ece01dccee9030181cddce8f35d0f95f82d70ff59d6e5ae78325f5a5cc9b9186be67efed24048cb6e6ad38605217242023a2745209e1ed7479cba400e6a2",
  "Transactions": null,
  "Nonce": "113805784024be8171a54b4eaef1a99d4d3cda00c256b157e9abbec85b58630eb6e1de04bafd4493309e3475519e0d7f08ca234faf33314a2264cb2ce56fd26e"
}
{
  "Hash": "0001c80d40bc6920c20cd8be6c224544833339221675ff0a4f76bd36cb14b308cf8c2aee62610d95b826298adecfde40c33b241b75c05f24a4209786eeac0e12",
  "Index": 32,
  "Timestamp": 1529426970,
  "Difficulty": 15,
  "CumulativeDifficulty": 256,
  "PrevHash": "00001d3c038a50185c39d46cbadc3da102779f4a8065a776cbf858b44483a5532095ff88c291bf154a02a6df0ac34c210ca931fa1c635c072ac7b66db06267bd",
  "Transactions": null,
  "Nonce": "f2327de645be39570389ca612f77008d5a77ed580e0c9d212995777d10ee10538060ce62dde2e867bf11428e4b862b87c9af29a33dfa1cca4e8a10130624878e"
}
{
  "Hash": "0000565302080f3521ace64a1328bd90a7709ac573da6f4dee4ad46e120d5f9c4f29bc9498eae7ec92a0ab2a15c71e9e7558872c1b8afad7aa6cbf02734eb201",
  "Index": 33,
  "Timestamp": 1529426973,
  "Difficulty": 17,
  "CumulativeDifficulty": 273,
  "PrevHash": "0001c80d40bc6920c20cd8be6c224544833339221675ff0a4f76bd36cb14b308cf8c2aee62610d95b826298adecfde40c33b241b75c05f24a4209786eeac0e12",
  "Transactions": null,
  "Nonce": "be07c7696318e05566ca2ad8a6cffccc618a3060dd0d4d0b841b86f2d3e6f171a382da3ee237d23486eefcb291c2980d64522081379efdadaa6462bf8b9699ea"
}
{
  "Hash": "0000162e55cdde9a79e39ad28ecfd401d81fe888264773f03718f83500ab7e7b7b8225c6c3168e6dbb76abf77bafebd38bd08a3c91032e6539d35230f318c6d6",
  "Index": 34,
  "Timestamp": 1529426973,
  "Difficulty": 16,
  "CumulativeDifficulty": 289,
  "PrevHash": "0000565302080f3521ace64a1328bd90a7709ac573da6f4dee4ad46e120d5f9c4f29bc9498eae7ec92a0ab2a15c71e9e7558872c1b8afad7aa6cbf02734eb201",
  "Transactions": null,
  "Nonce": "33dfa97f9af449fdd22b5630ffb3745167b522cbb7c902b4680f97366eba8a61355c9c8b0a63946d266901125ef0a2cca2bbfc18029b31a436e2d7bb7aac1537"
}
`

func CheckBlockLink(bc *Blockchain, block *Block, prev *Block) error {
	fmt.Printf("Index: %d Difficulty: %d\n", block.Index, block.Difficulty)
	ccdiff, _ := bc.GetCumulativeDifficulty(int(block.Index))
	pcdiff, _ := bc.GetCumulativeDifficulty(int(prev.Index))
	diffcmp := ccdiff == (pcdiff + int(block.Difficulty))
	if diffcmp == false {
		fmt.Printf("BAD CUMULATIVE DIFF: %d != %d + %d\n",
			ccdiff,
			pcdiff,
			block.Difficulty)
		return fmt.Errorf("bad difference")
	}
	return nil
}

func TestBlockchainAdd(t *testing.T) {
	bc := &Blockchain{}
	reader := bytes.NewReader([]byte(chain))
	decoder := json.NewDecoder(reader)

	block := &Block{}
	prev := &Block{}

	decoder.Decode(block)
	bc.Genesis(block)

	for {
		prev = block
		block = &Block{}
		err := decoder.Decode(&block)
		//fmt.Printf("BLOCKS=\n%+v\n%+v\n", prev, block)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
		}
		err = bc.Add(block)
		if err != nil {
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
			break
		}
		err = CheckBlockLink(bc, block, prev)
		if err != nil {
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
			break
		}
	}
}

func decodeBlocks(decoder *json.Decoder, count int) ([]*Block, error) {
	var blocks []*Block
	for i := 0; i < count; i++ {
		block := &Block{}
		if err := decoder.Decode(&block); err != nil {
			if err == io.EOF {
				break
			}
			return blocks, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func getBlockFragment(decoder *json.Decoder, count int) (*Fragment, error) {
	blocks, err := decodeBlocks(decoder, count)
	if err != nil {
		return nil, err
	}
	fragment := Fragment{
		Start:  blocks[0].Index,
		End:    blocks[len(blocks)-1].Index,
		Blocks: map[uint64]*Block{},
	}
	for _, block := range blocks {
		fragment.Blocks[block.Index] = block
	}
	return &fragment, nil
}

func TestBlockchainReplace(t *testing.T) {
	bc := &Blockchain{}
	reader := bytes.NewReader([]byte(chain))
	decoder := json.NewDecoder(reader)

	block := &Block{}

	decoder.Decode(block)
	bc.Genesis(block)

	f1c5, _ := getBlockFragment(decoder, 5)
	if err := bc.Replace(f1c5); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		t.Fail()
	}
	for i := 1; i < 5; i++ {
		prev := bc.Blocks[i-1]
		block := bc.Blocks[i]
		err := CheckBlockLink(bc, block, prev)
		if err != nil {
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
			break
		}
	}
	// Add a short diff forked block here ...
	f2c15, _ := getBlockFragment(decoder, 15)
	if err := bc.Replace(f2c15); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		t.Fail()
	}
	for i := 5; i < 15; i++ {
		prev := bc.Blocks[i-1]
		block := bc.Blocks[i]
		err := CheckBlockLink(bc, block, prev)
		if err != nil {
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
			break
		}
	}
}
