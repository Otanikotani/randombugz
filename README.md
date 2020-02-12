# Description
Simple program to generate cases in a given FogBugz instance.

# Running via Go
```
go run main.go --file-input <filepath> --fb-token <your-token> --fb-instance <your-fb-instance> 
```

Where:

`<your-token>` - is a API token that can be generated through web interface in FB. Go to User Options -> Create API Token

`<your-fb-instance>` - instance of the FB to insert cases into. For example https://sandbox.fogbugz.com

`<filepath>` - path to the .csv file with the following structure:

```csv
original;elapsed;milestone;user
1;1;1;1
1;1;1;1
1;1;1;1
1;1;1;1
1;1;1;1
1;1;1;1
1;1;1;1
1;1;1;1
1;1;1;2
1;1;1;2
1;1;1;2
1;1;1;2
1;1;1;2
1;1;1;2
1;1;1;2
1;1;1;2
8;0;1;1
8;0;1;1
8;0;1;2
8;0;1;2
``` 

This input will create:
- two new users
- one new project
- one new milestone
- new closed and estimated 1:1 original/elasped cases for both users
- new open and estimated 8 original cases for both users
 
Each line represents a case. If elapsed is zero - a case is created but not resolved/closed. Otherwise, the case will be resolved and closed.

Elapsed cannot be greater than 24
