
3.0.0Û
Test API for GSoC projectŒThis is a OpenAPI description for testing my GSoC project. The name of the path defines what
will be tested and the operation object will be set accordingly.
Structure of tests:
/testParameter*   --> To test everything related to path/query parameteres
/testResponse*    --> To test everything related to respones
/testRequestBody* --> To test everything related to request bodies
others            --> Other stuff

#TODO: ADD TESTS FOR components/requestBodies
21.0.0"é
ã
/testRequestBodyReferencen"l*testRequestBodyReference::8
6#/components/requestBodies/ComponentExampleRequestBodyB
200
	
success*«
ﬁ
€
ComponentExampleObjectPerson∫
∑*
Pet∫name∫	photoUrls object˙ê

id
 integeröint64

age
 integeröint64

name
:	doggie
 string
5
	photoUrls(
&*
photoUrl( arrayÚ

	 stringF
D
ComponentExampleParameter'
%
param1queryR
 integeröint64*õ
ò
ComponentExampleRequestBodyy
w
$A JSON object containing informationM
K
application/json7
53
1#/components/schemas/ComponentExampleObjectPerson