
3.0.0�
Test API for GSoC project�This is a OpenAPI description for testing my GSoC project. The name of the path defines what
will be tested and the operation object will be set accordingly.
Structure of tests:
/testParameter*   --> To test everything related to path/query parameteres
/testResponse*    --> To test everyting related to respones
others            --> Other stuff

#TODO: ADD TESTS FOR components/requestBodies
21.0.0"i
g
/testResponseNativeP"N*testResponseNativeB86
200/
-
succes#
!
application/json

	�string*�
�
�
ComponentExampleObjectPerson�
�*
Pet�name�	photoUrls�object��

id
�integer�int64

age
�integer�int64

name
:	doggie
�string
5
	photoUrls(
&*
photoUrl(�array�

	�stringF
D
ComponentExampleParameter'
%
param1queryR
�integer�int64