for i in tools peer orderer ccenv baseos ca baseimage
do 
	docker pull miguelgayo/fabric-$i	
	docker tag miguelgayo/fabric-$i hyperledger/fabric-$i
	docker tag miguelgayo/fabric-$i hyperledger/fabric-$i:2.1
done