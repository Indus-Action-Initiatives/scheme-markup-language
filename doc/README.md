# Scheme Markup Language (SML)
## Why?
SML is needed to manage the documentation and criteria of the schemes. A text file is the abstraction chosen to encapsulate the details about a scheme. It allows for version control, readability, easy extensions and standard software engineering principles like peer reviews, roll back, testing before productionising etc.

We have chosen to represent the SML as a indentation-sensitive key-value representation.

## Structure
Following are the components of a standard SML setup:

1. Project
2. Tables
3. Dataset
4. Schemes

### Project
A project is the outermost container holding all the logic and documentation regarding any scheme eligibility. It contains table files, dataset files and schemes files. There can be multiple table files and scheme files. But there has to be exactly one dataset file in the project. A project can be backed by a git repository to allow for version control and peer reviews, environments etc. You can imagine the project as the root directory in a filesystem.

### Table
Table files capture semantics of the logical entities from the beneficiary database. Usually it has one-to-one mapping with the beneficiary table. But depending on the structure of the schemes and beneficiary information structure, there can be multiple logical table. In this case, there will be multiple table files in the project e.g. if certain schemes need family information as well as the individual information for the same criteria, then family and the beneficiary will be a separate table; an example of such a criteria can be `family_income < 100000 AND beneficiary_gender = 'Male'`

### Dataset
A dataset includes entries for all the table files and relationship between them using **JOIN** conditions. There must be **exactly one** dataset file in the project.

### Scheme
A scheme has two main elements to it:

1. **Set of Criteria**: A criteria represents a particular evalutation evaluation criteria for a scheme. The criteria can have nested conditions.
	For example:
	
	> The age for males should be greater than or equal to 21 and for females it should be greater than equal to 18
	
	can be represented by
	
		criteria age
	        combine OR
	            combine AND
	                column gender
	                    operator equals
	                    value male
	                column age
	                    operator gte
	                    value 21
	            combine AND
	                column gender
	                    operator equals
	                    value female
	                column age
	                    operator gte
	                    value 18


2. **Evaluation**: An evalutation is a combination of the criteria in the scheme file using logical operators `&&` and `||`. Generally, the evaluation will be a conjunction of all the criteria in scheme file. There must be **exactly one** evaluation in a scheme file.

### Formatting columns
The numeric and date columns can be formatted using the standard Excel-like format strings.
For example:

> The percentage marks in the tenth & the twelfth class can be formatted as below:

	column tenth_percentage_marks
        type float
        sql tenthPercentageMarks
        format %.2f
    column twelfth_percentage_marks
        type float
        sql twelfth_percentage_marks
        format %.2f

### Auxiliary columns
Sometimes the values needed to evaluate the criteria do not come from the same row. They may come from a different row from the same time. In such cases, we need to generate and execute a query involving a *common table expression*.
For example:

> Imagine a table named family_members with each member as separate rows, with a column called family_role denoting the role of the member (father, mother, child, etc.) in the family. Suppose a criteria involves ascertaining if the father or mother of the child has a BOCW card. This criteria can be represented as follows:

	criteria parent_bocw
        combine OR            
            column mother_bocw
                table fm
                operator equals
                value True
            column father_bocw
                table fm
                operator equals
                value True
    criteria family_role
        column family_role
        table fm
        operator equals
        value child
> In the example above *mother_bocw* and *father_bocw* are auxiliary columns. Their definitions are as follows:

	auxiliary_column mother_bocw
        type bool
        sql SELECT mother.has_bocw_card FROM family_members as mother INNER JOIN families ON mother.family_id = families.id WHERE mother.family_role='mother' GROUP BY 1 LIMIT 1
    auxiliary_column father_bocw
        type bool
        sql SELECT father.has_bocw_card FROM family_members as father INNER JOIN families ON  father.family_id = fm.family_id WHERE father.family_role='father' GROUP BY 1 LIMIT 1

### Lambda transformers
Sometimes, SQL transformations aren't enough inorder to calculate the values needed to evaluate the criteria. These cases are the ones which are better solved using the more powerful languages like Golang. In EE, such transformations are called *Lambda transformations*.
For example:

> Collecting the pregnancy data can be quite tricky during campaigns. In order to get around the potentially offending question about the pregnancy status of a girl child, the survey is designed to collect the names of the pregnant members in the family if any. But now to attribute the *pregnancy_status* to the potential pregnant family member, one must deduce the precisely figure out which is the member of the family that is pregnant using the names collected during the survey.
> 
> This can be using the lambda transformers like the ones below. Following the definition of the column in the table file.

	column pregnancy_status
        type status
        transformer lambda
        name populatePregnancyStatus

> The corresponding lambda function for the transformer (in Python) is shown as follows:

	def populatePregnacyStatus(self, family, beneficiary):
        # get the names of the pregnant women of the family from the pregnancy mapping
        pregnantWomenCombinedDict = getMappedDict(
            config["pregnancyMapping"], beneficiary
        )
        # current data has a provision of only two pregnant women, so pass two
        pregnantWomen = splitCombinedDict(pregnantWomenCombinedDict, 2)

        for woman in pregnantWomen:
            if "name" not in woman:
                return
            name = woman["name"].strip()
            if name == "":
                return
            fuzzyScore, index = fuzzy_matching(name, family["members"])
            if index >= 0 and index < (len(family["members"]) - 1):
                family["members"][index]["pregnancy"] = "yes"
            # elif index == (len(family['members']) - 1):
            #     respondent['pregnancy'] = 'yes'
        for member in family["members"]:
            if member["gender"] == "male":
                member["pregnancy"] = "no"
            elif "pregnancy" not in member or member["pregnancy"] == "":
                member["pregnancy"] = UNKNOWN_STRING

# MakeSQL

MakeSQL deals with turning the project into a set of valid SQL statements that can be run against the in-memory DuckDB database, i.e. the sql queries are generated in the DuckDB dialect.

It has two main components:

### 1. Parser
The SML "parser" takes in the project in SML form and outputs a structured `JSON` representation of the project. It does more than just parsing including lexical analysis and semantic analysis as well. The SML parser thus also finds the inconsistencies in the project and flags them.

### 2. SQL Generator
The SQL generator takes in a "parsed" SML project in JSON format, a selection JSON (set of columns, filter strings etc.) and generates valid evaluation SQL to be executed against a populated DuckDB instance.

*Please note that itt is not a responsibility of MakeSQL to ensure that the SQL statements specified in the criteria are in the DuckDB namespace. MakeSQL assumes that it is the case.*


## SAMPLE input.json
TODO
### Ensuring that criteria SQL are in the DuckDB namespace
TODO

## DDL Queries
TODO

# Auto-generating a SML project from a sample input
TODO

