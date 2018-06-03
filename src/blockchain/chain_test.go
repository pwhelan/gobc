package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"testing"
)

var chain = `
{
	"Hash": "75d91be3fe7f83e50e7036aba4383b601404734db31ffffbddc215873289c78427fff7ca67fd4c83e79793b8cbb6440ee46760affadf48fc16473ee3e6424bb5",
	"Index": 0,
	"Timestamp": 1527949914,
	"Difficulty": 0,
	"PrevHash": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	"Transactions": null,
	"Nonce": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
}
{
	"Hash": "d6ceb96a4996773cacbebe846bbaf3f596b99f6e0d4889e8993decdb1d2fdb896b42f4c3d67469f011191eecc9b72b81a8e6579999e97fed4e0202e272b12fe2",
	"Index": 1,
	"Timestamp": 1527985615,
	"Difficulty": 0,
	"PrevHash": "75d91be3fe7f83e50e7036aba4383b601404734db31ffffbddc215873289c78427fff7ca67fd4c83e79793b8cbb6440ee46760affadf48fc16473ee3e6424bb5",
	"Transactions": null,
	"Nonce": "4b750e9bb8d5249e5bdfa92b061e99afa3ed7d522347466b62adec5cf247af3a39d8527b6927e40d2e2d1fd5c32268779fae7a55624d0188eb08efdb23137fad"
}
{
	"Hash": "8a8ca34a3a8a902a5edf608258176232b2b768b2fa802f8558168d284a1e9ddc3188fce4b0152fd5b3abd4c6d3a500f7185183625d7921c65b1fbfe5d93cf1c8",
	"Index": 2,
	"Timestamp": 1527985615,
	"Difficulty": 0,
	"PrevHash": "d6ceb96a4996773cacbebe846bbaf3f596b99f6e0d4889e8993decdb1d2fdb896b42f4c3d67469f011191eecc9b72b81a8e6579999e97fed4e0202e272b12fe2",
	"Transactions": null,
	"Nonce": "1058d0b2ba53f293c92bfac7980bc5bd313b1208a5540b782d4731d1c241ecb0e62a690bf2a76a76e3d7101a033cc22b5e31f34a9c638cd1baaeefa2d13a0bdd"
}
{
	"Hash": "1cb85a5e3f6722502a3cae273e3b0e9bd87101b0a349c17c41f990e7972f2036c5918fc311148a1df6fead107bff650a9af6b874696b4cd24253fef1b27e5a65",
	"Index": 3,
	"Timestamp": 1527985615,
	"Difficulty": 1,
	"PrevHash": "8a8ca34a3a8a902a5edf608258176232b2b768b2fa802f8558168d284a1e9ddc3188fce4b0152fd5b3abd4c6d3a500f7185183625d7921c65b1fbfe5d93cf1c8",
	"Transactions": null,
	"Nonce": "3e0be4174b58cf96e735a4cb59c17f3417d28ca4c9fea6750e291d10143ac9d2f21c066e211a4d94c58316a8c039e88868311b560dbc3af4b5a752547e370324"
}
{
	"Hash": "13fffdf38f469112805487786ac23fe5fcfbccf5d59df5b568069b84626e937c4e797be41850f14eae885c386036b1714f2780c3f77ac9dc9728e2c467e743e2",
	"Index": 4,
	"Timestamp": 1527985615,
	"Difficulty": 2,
	"PrevHash": "1cb85a5e3f6722502a3cae273e3b0e9bd87101b0a349c17c41f990e7972f2036c5918fc311148a1df6fead107bff650a9af6b874696b4cd24253fef1b27e5a65",
	"Transactions": null,
	"Nonce": "8204e6ae13ad06bee67379f1bc7f22ce38f6a06dccb5b9533ee85df2ca516b9156025900186c0450ba33480a532627fbf16c604af4a2d1e36e4198527a86dd90"
}
{
	"Hash": "197efefd3f20fe70cb4c4f64c8162afa9bb7e2e7e9a3469f2d66dbd8ad08825951c0e89c084becf998d6acede6a0dbc1bbde695261fb1f9edfd77918428df961",
	"Index": 5,
	"Timestamp": 1527985615,
	"Difficulty": 3,
	"PrevHash": "13fffdf38f469112805487786ac23fe5fcfbccf5d59df5b568069b84626e937c4e797be41850f14eae885c386036b1714f2780c3f77ac9dc9728e2c467e743e2",
	"Transactions": null,
	"Nonce": "4b5ba723d71fc7ca8030f0d6c3cd73531d45ef29a6d153abe1249495a55cf2abdd643069b7181f4f4ade3778ac5e2e8415ece537e110728c7aad25df48ce5190"
}
{
	"Hash": "0fa9fcdfb832ec64b160de2b1a7057915c7abbf2d03a6173672869eb262249780093c27476d5d30eae333291946c56e52aa04d070b94faf1424c2b833a122672",
	"Index": 6,
	"Timestamp": 1527985615,
	"Difficulty": 4,
	"PrevHash": "197efefd3f20fe70cb4c4f64c8162afa9bb7e2e7e9a3469f2d66dbd8ad08825951c0e89c084becf998d6acede6a0dbc1bbde695261fb1f9edfd77918428df961",
	"Transactions": null,
	"Nonce": "b902de48f87a62e4732d50fc8e734f245a387fbac7ef650a593d340a779dbd1c2dcdbd93c4a4d8af7edd5f995b77d0ec87b95448ded0300ac7dd2d345454b7d9"
}
{
	"Hash": "000d96ee06a72c5e5e557c88ce149b5f84c4e35345fc6d53293590de0e1b00e0e6f23c04ae8b37785121e192db3007adda2237f1a4b7096a043298ae6a0cbe61",
	"Index": 7,
	"Timestamp": 1527985615,
	"Difficulty": 5,
	"PrevHash": "0fa9fcdfb832ec64b160de2b1a7057915c7abbf2d03a6173672869eb262249780093c27476d5d30eae333291946c56e52aa04d070b94faf1424c2b833a122672",
	"Transactions": null,
	"Nonce": "93296afc54426ad53b6c4fe09307b425e9fc6f24ae25a230332c19d9b0284d1cc9bb1b7f6f96e29e3a93fd1643d86793d511db47a0e78ead4a18a1e9fcc99c69"
}
{
	"Hash": "00a9076b49fe97d6fb8b2984b553a6c7acb9bff328bc5863f83490fd8df3bb6c943e2ab291aeef10a4faa2e5a7def0bc224371d6dbbc2a04ac2707d2153b51de",
	"Index": 8,
	"Timestamp": 1527985615,
	"Difficulty": 6,
	"PrevHash": "000d96ee06a72c5e5e557c88ce149b5f84c4e35345fc6d53293590de0e1b00e0e6f23c04ae8b37785121e192db3007adda2237f1a4b7096a043298ae6a0cbe61",
	"Transactions": null,
	"Nonce": "6b6260e832e8bdd117b2eed0ece2fd29e3f7d276845c022dd141b02d22daee61c49c5020071454e5341340e7297f84ae168f32de455181d213c11b9d843ce8c5"
}
{
	"Hash": "00f4c6d34e0a53624009b8e36a4a44c2eefb1be71ab5597b7b2e65a8c80381f0d0657d21f7a0d4dfb582e7b810a6da26fa604b4a3c1b39977677ca0c79e69299",
	"Index": 9,
	"Timestamp": 1527985615,
	"Difficulty": 7,
	"PrevHash": "00a9076b49fe97d6fb8b2984b553a6c7acb9bff328bc5863f83490fd8df3bb6c943e2ab291aeef10a4faa2e5a7def0bc224371d6dbbc2a04ac2707d2153b51de",
	"Transactions": null,
	"Nonce": "15b114084309e347b2f57798d432b0309e86b5e0bbe0d5bad457145c76962915b91a8088116e9c812122746cb0f1277b87fb02ff11f830324351a2cb520b3d94"
}
{
	"Hash": "00c41cce11668d9584195c26c5628e99542cfc3f96f46aa063924bb71ed72cadd864e1dcd82f1a68a3a732436d67433d1dab78d0afc4dc3e2c3687a58e54492b",
	"Index": 10,
	"Timestamp": 1527985615,
	"Difficulty": 8,
	"PrevHash": "00f4c6d34e0a53624009b8e36a4a44c2eefb1be71ab5597b7b2e65a8c80381f0d0657d21f7a0d4dfb582e7b810a6da26fa604b4a3c1b39977677ca0c79e69299",
	"Transactions": null,
	"Nonce": "e5e70017b5d28f0c12e38aa5c536e18b010892ccf82f180cd4c83e7b92bf4034843cb76966d919263bf6447c52f590e8fdff2dac3f8b1ce14ecbc89df36fa75f"
}
{
	"Hash": "007cc61f672da653f55da0a415f7516f49ad5c31a1f19972eda0f4a89e774c8b2cba338f3bb2e4a11061c32cb2a164c83df191f0bb33e48cbda9bd4e46976bb9",
	"Index": 11,
	"Timestamp": 1527985615,
	"Difficulty": 9,
	"PrevHash": "00c41cce11668d9584195c26c5628e99542cfc3f96f46aa063924bb71ed72cadd864e1dcd82f1a68a3a732436d67433d1dab78d0afc4dc3e2c3687a58e54492b",
	"Transactions": null,
	"Nonce": "43abd65c1181a48640e0dd325118bc23e4335e2fbc8ad35e18030f7bd5edf1576c86d192f7fd076932846f2f88d7f5ded0f0e27255048ff8165a79322c71dd8e"
}
{
	"Hash": "0031d352a74e1ce052b8de88f04dd885700c2f78056c0e9d8f44fdd5bd0c28d47848767122181b0aaffd73c06d7d5cab6423d847e64db2f6b734c5a9310be5ad",
	"Index": 12,
	"Timestamp": 1527985615,
	"Difficulty": 10,
	"PrevHash": "007cc61f672da653f55da0a415f7516f49ad5c31a1f19972eda0f4a89e774c8b2cba338f3bb2e4a11061c32cb2a164c83df191f0bb33e48cbda9bd4e46976bb9",
	"Transactions": null,
	"Nonce": "ea90c99713e603c333a7400237f810a7d372e14782d37769bfb22c9e2e56b3f5a293738da3a6bfba1114f64ba96fa5396050d2c1213e30a29f6d16034fb7bfee"
}
{
	"Hash": "001033462e8f49a15d855483c6c7486f0f08b3d073e1acc5bec6fef69138be8410c656e0325365a68c820e8801bb3ca0dab6de9bfb90f80c41dd7b610f23b07e",
	"Index": 13,
	"Timestamp": 1527985616,
	"Difficulty": 11,
	"PrevHash": "0031d352a74e1ce052b8de88f04dd885700c2f78056c0e9d8f44fdd5bd0c28d47848767122181b0aaffd73c06d7d5cab6423d847e64db2f6b734c5a9310be5ad",
	"Transactions": null,
	"Nonce": "13ea2d1b8c22e74667548e4a827624970440213a1765ec1d7b8dfdd02af0380f4fa8d2f82c4fe0afd2c0c55edaa732e2a09f3c16dc7c6c0df524e8018c77bfbc"
}
{
	"Hash": "0004b3ab67c4ead0c83f47d078e067d363c07de196c2f673dca660afdbd10487a19fe3d6e1b4eb03ea55752aab57c5d7287e5d6cea141c99189a3b419a386c46",
	"Index": 14,
	"Timestamp": 1527985616,
	"Difficulty": 12,
	"PrevHash": "001033462e8f49a15d855483c6c7486f0f08b3d073e1acc5bec6fef69138be8410c656e0325365a68c820e8801bb3ca0dab6de9bfb90f80c41dd7b610f23b07e",
	"Transactions": null,
	"Nonce": "0a707d2c6ae0d0d895a18fefe743a746850e7af7d8a1f0cc48c561dcc2aea9cffc6a4d2b0ee9db261241aa8d12c70472dc19f7d2ab6758c3ef8f73e396e9b717"
}
{
	"Hash": "00013bbc5946ad4c08014a90c4c9d78b11ae04f7af7feaa5b8480810f317918b149cfab60476b1afe5ecb2b458721fda7577b3bacd5bd4af3d2090851a9c4087",
	"Index": 15,
	"Timestamp": 1527985616,
	"Difficulty": 13,
	"PrevHash": "0004b3ab67c4ead0c83f47d078e067d363c07de196c2f673dca660afdbd10487a19fe3d6e1b4eb03ea55752aab57c5d7287e5d6cea141c99189a3b419a386c46",
	"Transactions": null,
	"Nonce": "9c1366e465914cbcec3aa71c43fc86b435ae33b1f1a982a895f04b5c1e1bb0e10ca194f5126d0cdfe36e31eb9a7fe3893acfdcced346e24105aaf60e69c5d75c"
}
{
	"Hash": "0001963924f7b83301ad219576a8dc0f3d74ca3e81685abd51beaf955de1fe093e4b677add510686680a51559a5ad12d1ccef0fbc0ff5c687c1056eb4178a792",
	"Index": 16,
	"Timestamp": 1527985616,
	"Difficulty": 14,
	"PrevHash": "00013bbc5946ad4c08014a90c4c9d78b11ae04f7af7feaa5b8480810f317918b149cfab60476b1afe5ecb2b458721fda7577b3bacd5bd4af3d2090851a9c4087",
	"Transactions": null,
	"Nonce": "e5351eef45d4fd22166d9636f37bc2f9d8e3293fd3d9a4b577df7b77d5300acb7bf818c080ca22a4f851b6e249f43941706e561a7cd7b4469a4ac5bdbf4ac8b1"
}
{
	"Hash": "0000fb7636199f8420f096b669befd81830fe165bf819fdbf22893f45211c80d9cc8d13208dcf7a924bda88f4ce612f545e1f01b1769401bfb01ead40817ee87",
	"Index": 17,
	"Timestamp": 1527985617,
	"Difficulty": 15,
	"PrevHash": "0001963924f7b83301ad219576a8dc0f3d74ca3e81685abd51beaf955de1fe093e4b677add510686680a51559a5ad12d1ccef0fbc0ff5c687c1056eb4178a792",
	"Transactions": null,
	"Nonce": "a2d672ba4235737a61057e8ea3bc5289464b8604e290308d72814617aa7210cde3f1a408af2c68ac7482bca47d4c683aed630411a20b4814b494cb0cb7ecd103"
}
{
	"Hash": "000060b85c0cb372594c386641a745d044d8c14f19ecc133ad57f1612c9aa7c988e688df08d2f918e2ddc8a91d51c3201e53c4570baa40b3427984ec5d9d9764",
	"Index": 18,
	"Timestamp": 1527985617,
	"Difficulty": 16,
	"PrevHash": "0000fb7636199f8420f096b669befd81830fe165bf819fdbf22893f45211c80d9cc8d13208dcf7a924bda88f4ce612f545e1f01b1769401bfb01ead40817ee87",
	"Transactions": null,
	"Nonce": "e0f833e834d1970726588499679035aed835cb3acfa2603d62994685d32ca9122f5e67244eeeafb8938e38d6a33e79ddbdcd10a7b7a12ccce282ec99f319aba4"
}
{
	"Hash": "00002240b68cd9d014ec86c0413acbdfba75b89d1403a4ab5225e2dad633496e7a410a40133003d609feae51196d457eb592179719858d400c10f04643cc5d46",
	"Index": 19,
	"Timestamp": 1527985618,
	"Difficulty": 17,
	"PrevHash": "000060b85c0cb372594c386641a745d044d8c14f19ecc133ad57f1612c9aa7c988e688df08d2f918e2ddc8a91d51c3201e53c4570baa40b3427984ec5d9d9764",
	"Transactions": null,
	"Nonce": "60379deb33a5034e3b2562509aebf9a82f60a36774d19ba58dbd6586c277e75af2790ddb7499967f5ec0ed4bce3a315b0990ba238367163502b51c41adbc1590"
}
{
	"Hash": "00002d33b2c3867b799d5337865960435695a93e1abfcc85e16b90bad280a1a59747ba907227184a8fe5a6bb83ffcee911cb5c14c87e0fb78cbd42c82bd3283e",
	"Index": 20,
	"Timestamp": 1527985621,
	"Difficulty": 18,
	"PrevHash": "00002240b68cd9d014ec86c0413acbdfba75b89d1403a4ab5225e2dad633496e7a410a40133003d609feae51196d457eb592179719858d400c10f04643cc5d46",
	"Transactions": null,
	"Nonce": "05a25632b167fc985d9a6604fe1f7d52c96001936a874aa527ef845290c30c6fdc27aeb90ea56415478d5b33b783a8db98778dd79b9790107e7e48a708dee916"
}
{
	"Hash": "0000122b3ac1daf596ee8e5609492f0658bb3f03d95446e418c7a38a3fff4ea579eca2317491ffe1f0439d8dad915b3305fc0c97b6eb19744a7257871cd9d05d",
	"Index": 21,
	"Timestamp": 1527985624,
	"Difficulty": 19,
	"PrevHash": "00002d33b2c3867b799d5337865960435695a93e1abfcc85e16b90bad280a1a59747ba907227184a8fe5a6bb83ffcee911cb5c14c87e0fb78cbd42c82bd3283e",
	"Transactions": null,
	"Nonce": "2d73948e0f869aa5d1170e2553c0fdedbcb3ccd4c88c70c744761fdb2e0d2f7e95436c87f1c56ffa881f0904f9eb1d32b9136427ed4ccdc38fcfd493732ee007"
}
{
	"Hash": "00000e859e8d909d77f2280855cdb8016e7ef45742da6a3d71ed6779491f85eca9da47e1cb44c55a5588ed84e44a4ca89d7d95617afe9ae2541c5ca3ae5a9ab3",
	"Index": 22,
	"Timestamp": 1527985657,
	"Difficulty": 20,
	"PrevHash": "0000122b3ac1daf596ee8e5609492f0658bb3f03d95446e418c7a38a3fff4ea579eca2317491ffe1f0439d8dad915b3305fc0c97b6eb19744a7257871cd9d05d",
	"Transactions": null,
	"Nonce": "ebd7ab175192b9d6a055cd75422a60822ac84eee31ddb054b39fa787a4cd654b13b861780db55bf1e433acde3331a08bb4ad9d946ea38c45bfadb9924a3d949a"
}
{
	"Hash": "000005740a30fd13f02e084ce17238c66186e56dcfd941c042f05816aed2a7707c6e199957188a3b859b5336422b7fb65621f5f4e353083de528aef7b5740f43",
	"Index": 23,
	"Timestamp": 1527985724,
	"Difficulty": 21,
	"PrevHash": "00000e859e8d909d77f2280855cdb8016e7ef45742da6a3d71ed6779491f85eca9da47e1cb44c55a5588ed84e44a4ca89d7d95617afe9ae2541c5ca3ae5a9ab3",
	"Transactions": null,
	"Nonce": "9340b0a891d37c9aaf5f85b3fbed1106479c9d88251a34a4ad3cd55c65813eb6062f04d8ae8048772546d40c94a70ba3268060c8cef2e317f47190750076a9be"
}
{
	"Hash": "0000022cd760a555fd3e2bc6b89d15497dc17312a01d95c293eb904f6e6c2955d274a1cf16c5a5b6b4b833b1b05f054a4497d931a864d65acb49c25a6ed1892b",
	"Index": 24,
	"Timestamp": 1527985733,
	"Difficulty": 22,
	"PrevHash": "000005740a30fd13f02e084ce17238c66186e56dcfd941c042f05816aed2a7707c6e199957188a3b859b5336422b7fb65621f5f4e353083de528aef7b5740f43",
	"Transactions": null,
	"Nonce": "62d4667a11928e8fd544add87d9b7f73b0f3c8e6c2f2d2bc952c9536926f04996f56086108171eb2c8bbf73987e629a9069170504c3c2d59714baf3908c77d44"
}
{
	"Hash": "000000414a6124baa6935ff2e97969919305d1647a9ab82fa8d0e93d4cecbd1a549f83301ecd27446f717d316f1c4e9c9ab9d46f85b02566555adfbff6ac4613",
	"Index": 25,
	"Timestamp": 1527985852,
	"Difficulty": 23,
	"PrevHash": "0000022cd760a555fd3e2bc6b89d15497dc17312a01d95c293eb904f6e6c2955d274a1cf16c5a5b6b4b833b1b05f054a4497d931a864d65acb49c25a6ed1892b",
	"Transactions": null,
	"Nonce": "82988bc718efd7c53e6448fe11bdff5d6daedb7aea71c47415751be1bbcd8f6fb153f34353f337cebecad2593eb8177c80c0851b70fbfb556972c8c3725c953c"
}
{
	"Hash": "00000078485e8cb7e19652a1f6fa038f297b54b1c3f393f0269af56bebdf30b0e77ca96617e28e30e35abeb7ea0ce1fee9e5059e286ea1673aa65ebb920bd420",
	"Index": 26,
	"Timestamp": 1527985942,
	"Difficulty": 23,
	"PrevHash": "000000414a6124baa6935ff2e97969919305d1647a9ab82fa8d0e93d4cecbd1a549f83301ecd27446f717d316f1c4e9c9ab9d46f85b02566555adfbff6ac4613",
	"Transactions": null,
	"Nonce": "25e1862dbb287f0e90965dde688d442ff0d5226bcc50d84e63eec80122d08ec0a29839a49039bc5aa19db42fe5c94251da5f97862e32fed1418efc43cea05b3f"
}
{
	"Hash": "000000b85ab6982a6fcd8190613ee8e84d107c66a3667e2412d8d263ec70912e69d287427a3dd91560dca64c62f319fb7bd5b22c71fcdbd4e23e83ab3058197d",
	"Index": 27,
	"Timestamp": 1527986020,
	"Difficulty": 24,
	"PrevHash": "00000078485e8cb7e19652a1f6fa038f297b54b1c3f393f0269af56bebdf30b0e77ca96617e28e30e35abeb7ea0ce1fee9e5059e286ea1673aa65ebb920bd420",
	"Transactions": null,
	"Nonce": "a8787b8e77d170dc63d46960145ab0bdde9f5f9d08609fc08fdfcd3e783a031b3df29b8443fe5495fc0a0018f2c6af64fce2005a555c2c7aab794e9da0bf43b2"
}
{
	"Hash": "00000037e0a9888b74f155dece251795545e1ca2a818c3d8c266d2cd726852728d6ce1640a30ed168dffc4de2bd7a5f22910dfbcbb99b01240772311cb88e460",
	"Index": 28,
	"Timestamp": 1527986056,
	"Difficulty": 25,
	"PrevHash": "000000b85ab6982a6fcd8190613ee8e84d107c66a3667e2412d8d263ec70912e69d287427a3dd91560dca64c62f319fb7bd5b22c71fcdbd4e23e83ab3058197d",
	"Transactions": null,
	"Nonce": "e03a3088ac2ee2c8dcb8edfc5adb68346191b8f77409d08cb4ebdfb7ba3e3b9561c359f204ff05d0a76255d83945fbf7de59936ecf00689cfa30545e92dcd7f4"
}
`

func CheckBlockLink(block *Block, prev *Block) error {
	fmt.Printf("Index: %d Difficulty: %d CumulativeDifficulty: %s PrevDiff: %s\n",
		block.Index, block.Difficulty, block.CumulativeDifficulty.String(),
		prev.CumulativeDifficulty)
	diffcmp := block.CumulativeDifficulty.Cmp(
		big.NewInt(0).Add(
			prev.CumulativeDifficulty,
			big.NewInt(int64(block.Difficulty),
		)))
	if diffcmp != 0 {
		fmt.Printf("BAD CUMULATIVE DIFF: %s != %s + %d\n",
			block.CumulativeDifficulty.String(),
			prev.CumulativeDifficulty.String(),
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
		err := decoder.Decode(&block);
		//fmt.Printf("BLOCKS=\n%+v\n%+v\n", prev, block)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
		}
		err = bc.Add(block);
		if err != nil {
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
			break
		}
		err = CheckBlockLink(block, prev)
		if err != nil {
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
			break
		}
	}
}

func decodeBlocks(decoder *json.Decoder, count int) ([]*Block, error) {
	var blocks  []*Block
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
		Start: blocks[0].Index,
		End: blocks[len(blocks)-1].Index,
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
		err := CheckBlockLink(block, prev)
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
		err := CheckBlockLink(block, prev)
		if err != nil {
			fmt.Printf("ERROR[%d]: %s\n", block.Index, err.Error())
			t.Fail()
			break
		}
	}
}
