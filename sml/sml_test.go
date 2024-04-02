package sml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestParse0(t *testing.T) {
	input := `
project blah
  node foo`
	w, e := ParseSML(input, "")
	if e == nil {
		t.Errorf("expected error, got: %v (line=%d)", w, e[0].LineNum)
	} else {
		const expectedError = "error at position 3: unknown keyword: node"
		if e[0].Msg != expectedError {
			t.Errorf("expected error msg: %s\ngot: %s", expectedError, e[0].Msg)
		}
		if 3 != e[0].LineNum {
			t.Errorf("expected errorLine: 3\ngot: %d", e[0].LineNum)
		}
	}
}

func TestParse1(t *testing.T) {
	input := `
project blah
	table foo
		sql bar`
	_, e := ParseSML(input, "")
	if e != nil {
		t.Errorf("got error: %v", e[0].Msg)
	}
	_, e = ParseSML(strings.Replace(input, "blah", "", 1), "")
	const expectedError = `"project" has to have a name`
	if e[0].Msg != expectedError {
		t.Errorf("expected error: %s; got: %s", expectedError, e[0].Msg)
	}
}

func TestParse2(t *testing.T) {
	input := `
project bocw
    dataset main
        table fm
    table fm
        description "table for beneficiary data"
		sql fm
        column id
            type string
            sql id
        column age
            type int
            sql age
`
	w, e := ParseSML(input, "")
	if e != nil {
		t.Errorf("Parse returned error: %s", e[0].Msg)
		return
	}
	b, e1 := json.Marshal(w)
	if e1 != nil {
		t.Errorf("json Marshal returned error: %v", e1)
		return
	}

	expectedJSON, _ := ioutil.ReadFile("testdata/sml02.json")
	diff := DiffJSON(string(expectedJSON), string(b))
	if diff != "" {
		fmt.Printf("%v\n", string(b))
		t.Errorf("obtained json is not expected; diff=%s", diff)
	}
}

func TestParse3(t *testing.T) {
	input := `
project bocw
    dataset main
        table fm
        table f
        table locations
        join fm <-> f
            on f.id = fm.family_id
        join f <-> locations
            on f.location_id = locations.id 
    table fm
        description "table for beneficiary data"
		sql fm
        column id
            type string
            sql id
        column age
            type int
            sql age
        column gender
            type string
            sql gender
        column pregnancy_status
            type string
            sql pregnancy_status
        column marital_status
            type string
            sql marital_status
        column bocw_card_registration_date
            type datetime
            sql bocw_card_registration_date
        column health_status
            type string
            sql health_status
        column number_of_children
            type int
            sql number_of_children
        column children_school_or_college
            type string
            sql children_school_or_college
        column spouse_alive
            type bool
            sql spouse_alive
        column occupation_of_surviving_spouse
            type string
            sql occupation_of_surviving_spouse
        column spouse_occupation // <---- This column name is used in the scheme criteria
            type string
            sql occupation_of_surviving_spouse
        column receiving_pension
            type bool
            sql receiving_pension
        column receiving_government_aid
            type bool
            sql receiving_government_aid
        column home_ownership_status
            type string
            sql home_ownership_status
	table f
		description "table for beneficiary family details"
		sql f
		column id
			type string
			sql id
		column location_id
			type string
			sql locationID        
		column caste
			type string
			sql caste
		column pr_of_cg
			type bool
			sql prOfCG
		column has_residence_certificate
			type bool
			sql hasResidenceCertificate
		column ration_card_type
			type string
			sql rationCardType
		column ptgo_or_pvtg
			type bool
			sql ptgoOrPVTG
		column are_forest_dwellers
			type bool
			sql areForestDwellers
		column has_phone
			type bool
			sql hasPhone
		column neighbourhood_phone
			type bool
			sql CASE WHEN neighbourhoodPhone <> "" THEN TRUE ELSE FALSE END
	table locations
		sql locations
		column id
			type string
			sql id
		column area_type
			type string
			sql area_type
		column pincode
			type string
			sql pincode
		column ward_number
			type string
			sql wardNumber
		column ward_name
			type string
			sql wardName
		column village
			type string
			sql village
		column surveyVillageTownCity
			type string
			sql surveyVillageTownCity
`
	w, e := ParseSML(input, "")
	if e != nil {
		t.Errorf("Parse returned error: %s", e[0].Msg)
		return
	}
	b, e1 := json.Marshal(w)
	if e1 != nil {
		t.Errorf("json Marshal returned error: %v", e1)
		return
	}

	expectedJSON, _ := ioutil.ReadFile("testdata/sml03.json")
	diff := DiffJSON(string(expectedJSON), string(b))
	if diff != "" {
		fmt.Printf("%v\n", string(b))
		t.Errorf("obtained json is not expected; diff=%s", diff)
	}
}

func TestParse4(t *testing.T) {
	input := `
project bocw
    dataset main
        table fm
        table f
        table locations
        join fm <-> f
            on f.id = fm.family_id
        join f <-> locations
            on f.location_id = locations.id 
    table fm
        description "table for beneficiary data"
		sql fm
        column id
            type string
            sql id
        column has_bocw_card
            type bool
            sql has_bocw_card
        column dob
            type datetime
            sql dob
        column age
            type int
            sql age
        column gender
            type string
            sql gender
        column pregnancy_status
            type string
            sql pregnancy_status
        column marital_status
            type string
            sql marital_status
        column bocw_card_registration_date
            type datetime
            sql bocw_card_registration_date
        column health_status
            type string
            sql health_status
        column number_of_children
            type int
            sql number_of_children
        column children_school_or_college
            type string
            sql children_school_or_college
        column spouse_alive
            type bool
            sql spouse_alive
        column occupation_of_surviving_spouse
            type string
            sql occupation_of_surviving_spouse
        column spouse_occupation // <---- This column name is used in the scheme criteria
            type string
            sql occupation_of_surviving_spouse
        column receiving_pension
            type bool
            sql receiving_pension
        column receiving_government_aid
            type bool
            sql receiving_government_aid
        column home_ownership_status
            type string
            sql home_ownership_status
	table f
		description "table for beneficiary family details"
		sql f
		column id
			type string
			sql id
		column location_id
			type string
			sql locationID        
		column caste
			type string
			sql caste
		column pr_of_cg
			type bool
			sql prOfCG
		column has_residence_certificate
			type bool
			sql hasResidenceCertificate
		column ration_card_type
			type string
			sql rationCardType
		column ptgo_or_pvtg
			type bool
			sql ptgoOrPVTG
		column are_forest_dwellers
			type bool
			sql areForestDwellers
		column has_phone
			type bool
			sql hasPhone
		column neighbourhood_phone
			type bool
			sql CASE WHEN neighbourhoodPhone <> "" THEN TRUE ELSE FALSE END
	table locations
		sql locations
		column id
			type string
			sql id
		column area_type
			type string
			sql area_type
		column pincode
			type string
			sql pincode
		column ward_number
			type string
			sql wardNumber
		column ward_name
			type string
			sql wardName
		column village
			type string
			sql village
		column surveyVillageTownCity
			type string
			sql surveyVillageTownCity
	scheme sewing_machine
		label Chief Minister Sewing Machine Assistance Scheme
		description """Female construction worker receives sewing machine"""    
		criteria bocw_card
			column has_bocw_card
			table fm
			operator equals
			value true
		criteria age       
			column dob
			table fm
			operator age_between
			value [18, 50]
			granularity year
		criteria gender
			column gender
			table fm
			operator equals
			value female
		
		evaluation bocw_card && age && gender
				
`
	w, e := ParseSML(input, "")
	if e != nil {
		t.Errorf("Parse returned error: %s", e[0].Msg)
		return
	}
	b, e1 := json.Marshal(w)
	if e1 != nil {
		t.Errorf("json Marshal returned error: %v", e1)
		return
	}

	expectedJSON, _ := ioutil.ReadFile("testdata/sml04.json")
	diff := DiffJSON(string(expectedJSON), string(b))
	if diff != "" {
		fmt.Printf("%v\n", string(b))
		t.Errorf("obtained json is not expected; diff=%s", diff)
	}
}

func TestParse5(t *testing.T) {
	input := `
project bocw
    dataset main
        table fm
        table f
        table locations
        join fm <-> f
            on f.id = fm.family_id
        join f <-> locations
            on f.location_id = locations.id 
    table fm
        description "table for beneficiary data"
		sql fm
        column id
            type string
            sql id
        column age
            type int
            sql age
        column dob
            type datetime
            sql dob
        column gender
            type string
            sql gender
        column pregnancy_status
            type string
            sql pregnancy_status
        column marital_status
            type string
            sql marital_status
        column has_bocw_card
            type bool
            sql has_bocw_card
        column bocw_card_registration_date
            type datetime
            sql bocw_card_registration_date
        column health_status
            type string
            sql health_status
        column number_of_children
            type int
            sql number_of_children
        column children_school_or_college
            type string
            sql children_school_or_college
        column spouse_alive
            type bool
            sql spouse_alive
        column father_bocw
            type bool
            sql father_bocw
        column occupation_of_surviving_spouse
            type string
            sql occupation_of_surviving_spouse
        column spouse_occupation // <---- This column name is used in the scheme criteria
            type string
            sql occupation_of_surviving_spouse
        column receiving_pension
            type bool
            sql receiving_pension
        column receiving_government_aid
            type bool
            sql receiving_government_aid
        column home_ownership_status
            type string
            sql home_ownership_status
	table f
		description "table for beneficiary family details"
		sql f
		column id
			type string
			sql id
		column location_id
			type string
			sql locationID        
		column caste
			type string
			sql caste
		column pr_of_cg
			type bool
			sql prOfCG
		column has_residence_certificate
			type bool
			sql hasResidenceCertificate
		column ration_card_type
			type string
			sql rationCardType
		column ptgo_or_pvtg
			type bool
			sql ptgoOrPVTG
		column are_forest_dwellers
			type bool
			sql areForestDwellers
		column has_phone
			type bool
			sql hasPhone
		column neighbourhood_phone
			type bool
			sql CASE WHEN neighbourhoodPhone <> "" THEN TRUE ELSE FALSE END
	table locations
		sql locations
		column id
			type string
			sql id
		column area_type
			type string
			sql area_type
		column pincode
			type string
			sql pincode
		column ward_number
			type string
			sql wardNumber
		column ward_name
			type string
			sql wardName
		column village
			type string
			sql village
		column surveyVillageTownCity
			type string
			sql surveyVillageTownCity
	scheme sewing_machine
		label Chief Minister Sewing Machine Assistance Scheme
		description """Female construction worker receives sewing machine"""
		criteria bocw_card
			column has_bocw_card
			table fm
			operator equals
			value true
		criteria age       
			column dob
			table fm
			operator age_between
			value [18, 50]
			granularity year
		criteria gender
			column gender
			table fm
			operator equals
			value female
		
		evaluation bocw_card && age && gender
	scheme pm_ujjwala
		label Pradhan Mantri Ujjwala Yojana
		description """Free gas cylinders and stoves are provided to registered women laborers or registered construction laborers through the Food Department to the laborer's wife."""        
		criteria bocw
			combine OR
				term bocw_card
					column has_bocw_card
					table fm
					operator equals
					value True
				term father_bocw
					column father_bocw
					table fm
					operator equals
					value True
		criteria gender
			column gender
			table fm
			operator equals
			value female
		
		evaluation gender && bocw
				
`
	w, e := ParseSML(input, "")
	if e != nil {
		t.Errorf("Parse returned error: %s", e[0].Msg)
		return
	}
	b, e1 := json.Marshal(w)
	if e1 != nil {
		t.Errorf("json Marshal returned error: %v", e1)
		return
	}

	expectedJSON, _ := ioutil.ReadFile("testdata/sml05.json")
	diff := DiffJSON(string(expectedJSON), string(b))
	if diff != "" {
		fmt.Printf("%v\n", string(b))
		t.Errorf("obtained json is not expected; diff=%s", diff)
	}
}

func TestParseBoCW(t *testing.T) {
	input := `
project bocw
    dataset main
        table fm
    table fm
        label fm
		sql fm
        description "table for beneficiary data"
        column id
            type string
            sql id
        column age
            type int
            sql age
        column gender
            type string
            sql gender
        column pregnancy_status
            type string
            sql pregnancy_status
        column occupation
            type string
            sql occupation
        column marital_status
            type string
            sql marital_status
        column bocw_card_registration_date
            type datetime
            sql bocw_card_registration_date
        column health_status
            type string
            sql health_status
        column number_of_children
            type int
            sql number_of_children
        column children_school_or_college
            type string
            sql children_school_or_college
        column spouse_alive
            type bool
            sql spouse_alive
        column occupation_of_surviving_spouse
            type string
            sql occupation_of_surviving_spouse
        column spouse_occupation // <---- This column name is used in the scheme criteria
            type string
            sql occupation_of_surviving_spouse
        column receiving_pension
            type bool
            sql receiving_pension
        column receiving_government_aid
            type bool
            sql receiving_government_aid
        column home_ownership_status
            type string
            sql home_ownership_status
    scheme death_benefit
        label Compensation in case of Death (Accident/Natural Death)
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria health_status
            table fm
            column health_status
            operator equals
            value Deceased
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 0
    
        evaluation occupation && health_status && card_registration
    scheme disability_compensation
        label Compensation for disability
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria health_status
            table fm
            column health_status
            operator equals
            value Permanent disability as per disability certificate
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 0
    
        evaluation occupation && health_status && card_registration
    scheme disability_pension
        label Disability Pension 
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria health_status
            table fm
            column health_status
            operator equals
            value Permanent disability as per disability certificate
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 0
    
        evaluation occupation && health_status && card_registration
    scheme education_assistance_school
        label Education Assistance (School)
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria number_of_children
            table fm
            column number_of_children
            operator gte
            value 1
        criteria children_in_college
            table fm        
            column children_school_or_college
            operator equals
            value college
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 0
    
        evaluation occupation && number_of_children && children_in_college && card_registration
    scheme education_assistance_school
        label Education Assistance (School)
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria number_of_children
            table fm
            column number_of_children
            operator gte
            value 1
        criteria children_in_school
            table fm
            column children_school_or_college
            operator equals
            value school
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 0
    
        evaluation occupation && number_of_children && children_in_school && card_registration
    scheme family_pension
        label Family Pension
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria health_status
            table fm
            column health_status
            operator equals
            value Deceased
        criteria spouse_alive
            table fm
            column spouse_alive
            operator equals
            value True
        criteria spouse_occupation
            table fm
            column spouse_occupation
            operator equals
            value Pensioner
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 1
    
        evaluation occupation && health_status && spouse_alive && spouse_occupation && card_registration
    scheme family_pension
        label Family Pension
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria health_status
            table fm
            column health_status
            operator equals
            value Deceased
        criteria spouse_alive
            table fm
            column spouse_alive
            operator equals
            value True
        criteria spouse_occupation
            table fm
            column spouse_occupation
            operator equals
            value Pensioner
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 1

        evaluation occupation && health_status && spouse_alive && spouse_occupation && card_registration
    scheme funeral_assistance
        label Funeral Assistance
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria health_status
            table fm
            column health_status
            operator equals
            value Deceased
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 0

        evaluation occupation && health_status && card_registration
    scheme house_purchase_advance
        label Advance for purchase or construction of the house
        criteria age
            table fm
            column age
            operator lt
            value 50
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria home_ownership_status
            table fm
            column home_ownership_status
            operator IN
            value ['Owner', 'Going to purchase']
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 5

        evaluation age && occupation && home_ownership_status && card_registration
    scheme maternity_benefit
        label Maternity Benefit
        description """Maternity benefits of Rs. 30,000 to registered women members 
    and wives of male members (upto 2 children). (Rule – 271) – from the date of joining membership of the fund."""
        criteria gender_and_age
            combine OR
                term male_gender_and_age
                    combine AND
                        term male_gender
                            table fm
                            column gender
                            operator equals
                            value male
                        term male_age
                            table fm
                            column age
                            operator gte
                            value 21   
                term female_gender_and_age         
                    combine AND
                        term female_gender
                            table fm
                            column gender
                            operator equals
                            value female
                        term female_age
                            table fm
                            column age
                            operator gte
                            value 18
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria marital_status
            table fm
            column marital_status
            operator equals
            value Married
        criteria pregnancy_status
            table fm
            column pregnancy_status
            operator IN
            value ['Delivered first child', 'Delivered second child']
        criteria number_of_children
            table fm
            column number_of_children
            operator gte
            value 1        
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 0

        evaluation age && occupation && marital_status && pregnancy_status && number_of_children && card_registration
    scheme medical_assistance
        label Medical Assistance
        description """Financial Assistance for marriage of self and for children (upto 2 children). (Rule – 282), the building workers having continuous membership of 03 years shall be eligible. The details are as under:-
    Marriage of female registered member – Rs.51,000/-
    Marriage of male registered member    - Rs.35,000/-
    Marriage of daughter of registered members – Rs.51,000/-
    Marriage of son of registered members – Rs.35,000/
    and wives of male members (upto 2 children). (Rule – 271) – from the date of joining membership of the fund."""    
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria health_status
            table fm
            column health_status
            operator IN
            value ['Hospitalised for more than 5 days', 'In plaster at residence']
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 0

        evaluation occupation && health_status && card_registration
    scheme miscarriage_compensation
        label Compensation in case of miscarriage
        criteria gender_and_age
            combine OR
                term male_gender_and_age
                    combine AND
                        term male_gender
                            table fm
                            column gender
                            operator equals
                            value male
                        term male_age
                            table fm
                            column age
                            operator gte
                            value 21   
                term female_gender_and_age         
                    combine AND
                        term female_gender
                            table fm
                            column gender
                            operator equals
                            value female
                        term female_age
                            table fm
                            column age
                            operator gte
                            value 18
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria marital_status
            table fm
            column marital_status
            operator equals
            value Married
        criteria pregnancy_status
            table fm
            column pregnancy_status
            operator IN
            value ['Miscarriage', 'Still born']
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte
            granularity year
            value 0

        evaluation age && occupation && marital_status && pregnancy_status && card_registration
    scheme old_age_pension
        label Old Age Pension Benefit
        criteria age
            table fm
            column age
            operator gt
            value 60
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria health_status
            table fm
            column health_status
            operator ne
            value Deceased
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 3

        evaluation age && occupation && health_status && card_registration
    scheme work_tools_loan
        label Loan for purchase of work tools
        criteria age
            table fm
            column age
            operator lt
            value 55
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 3

        evaluation age && occupation && card_registration
    scheme work_tools_grant
        label Grant for purchase of work-related tools
        criteria age
            table fm
            column age
            operator lt
            value 55
        criteria occupation
            table fm
            column occupation
            operator equals
            value Construction Worker
        criteria card_registration
            table fm
            column bocw_card_registration_date
            operator age_gte // generalising age to be difference between the column and today
            granularity year
            value 3

        evaluation age && occupation && card_registration
`
	w, e := ParseSML(input, "")
	if e != nil {
		t.Errorf("Parse returned error: %s", e[0].Msg)
		return
	}
	b, e1 := json.Marshal(w)
	if e1 != nil {
		t.Errorf("json Marshal returned error: %v", e1)
		return
	}

	expectedJSON, _ := ioutil.ReadFile("testdata/smlBoCW.json")
	diff := DiffJSON(string(expectedJSON), string(b))
	if diff != "" {
		fmt.Printf("%v\n", string(b))
		t.Errorf("obtained json is not expected; diff=%s", diff)
	}
}

func TestParseCGRTEPlus(t *testing.T) {
	input := `
project cg_rte_plus
	dataset main
		table fm
		table f
		table locations
		join fm <-> f        
			on f.id = fm.family_id
		join f <-> locations
			on f.location_id = locations.id 
    table fm
		sql fm
        label fm
        description "table for beneficiary data"
        column id
            type string
            sql id
        column family_id
            type string
            sql familyID        
        column name
            type string
            sql name
        column dob
            type datetime
            format "%Y-%m-%dT%H:%M:%S.%fZ"
            sql dob
		column family_rank
			type string
			sql family_rank
        column gender
            type string
            sql gender
        column family_role
            type string
            transformer lambda
            transformer_name findFamilyRole
        column relationship_with_respondent
            type string
            sql relationshipWithRespondent    
        column disadvantaged
            type bool
            sql disadvantaged
        column pregnancy_status
            type string
            transformer lambda
            transformer_name populatePregnancyStatus
        column job
            type string
            sql job
        column job_type
            type string
            sql jobType
        column in_educational_institute
            type bool
            sql inEducationalInstitute
        column education_level
            type string
            sql educationLevel
        column prev_year_tenth
            type bool
            sql prevYearTenth        
        column prev_year_twelfth
            type bool
            sql prevYearTwelfth
        column tenth_percentage_marks
            type float
            sql tenthPercentageMarks
            format %.2f
        column twelfth_percentage_marks
            type float
            sql twelfth_percentage_marks
            format %.2f
        column tenth_top_ten
            type bool
            sql tenth_top_ten
        column twelfth_top_ten
            type bool
            sql twelfthTopTen
        column has_bocw_card
            type bool
            sql hasBOCWCard
        column bocw_card_issue_date
            type datetime
            format "%Y-%m-%dT%H:%M:%S.%fZ"
            sql bocwCardIssueDate
        column has_uow_card
            type bool
            sql hasUOWCard
        column uow_card_issue_date
            type datetime
            format "%Y-%m-%dT%H:%M:%S.%fZ"
            sql uowCardIssueDate

        // auxiliary columns start here
        auxiliary_column mother_bocw
            type bool
            sql SELECT mother.has_bocw_card FROM family_members as mother INNER JOIN families ON mother.family_id = families.id WHERE mother.family_role='mother' GROUP BY 1 LIMIT 1
        auxiliary_column father_bocw
            type bool
            sql SELECT father.has_bocw_card FROM family_members as father INNER JOIN families ON  father.family_id = fm.family_id WHERE father.family_role='father' GROUP BY 1 LIMIT 1    
    table f
        label f
		sql f
        description "table for beneficiary family details"
        column id
            type string
            sql id            
        column location_id
            type string
            sql locationID        
        column caste
            type string
            sql caste
        column pr_of_cg
            type bool
            sql prOfCG
        column has_residence_certificate
            type bool
            sql hasResidenceCertificate
        column ration_card_type
            type string
            sql rationCardType
        column ptgo_or_pvtg
            type bool
            sql ptgoOrPVTG
        column are_forest_dwellers
            type bool
            sql areForestDwellers
        column has_phone
            type bool
            sql hasPhone
        column neighbourhood_phone
            type bool
            sql CASE WHEN neighbourhoodPhone <> "" THEN TRUE ELSE FALSE END
    table locations
		sql locations
        label locations
        column id
            type string
            sql id
        column area_type
            type string
            sql area_type
        column pincode
            type string
            sql pincode
        column ward_number
            type string
            sql wardNumber
        column ward_name
            type string
            sql wardName
        column village
            type string
            sql village
        column surveyVillageTownCity
            type string
            sql surveyVillageTownCity
    scheme cycle_assistance
        label Chief Minister Cycle Assistance Scheme
        description """Construction worker receives cycle"""    
        criteria bocw_card        
            column has_bocw_card
            table fm
            operator equals
            value true
        criteria age_by_gender
            combine OR
                term female
                    combine AND
                        term female_age
                            column dob
                            table fm
                            operator age_between
                            value [18, 35]
                            granularity year
                        term female_gender
                            column gender
                            table fm
                            operator equals
                            value female
                term male
                    combine AND
                        term male_age
                            column dob
                            table fm
                            operator age_between
                            value [18, 50]
                            granularity year
                        term male_gender
                            column gender
                            table fm
                            operator equals
                            value male
        
        evaluation bocw_card && age_by_gender
    scheme merit
        label Meritorious Student / Student Education Promotion Scheme
        description """Children of construction workers get INR 2,000 to INR 12,500 if they perform well in class 10th or 12th Chhattisgarh Board exams. An additional benefit of 1,00,000 is given if the child is in the top 10 of merit list."""
        criteria parent_bocw
            combine OR
                term mother_bocw      
                    column mother_bocw
                    table fm
                    operator equals
                    value True
                term father_bocw
                    column father_bocw
                    table fm
                    operator equals
                    value True
        criteria family_role
            column family_role
            table fm
            operator equals
            value child
        criteria merit
            combine OR
                term tenth_merit
                    column tenth_percentage_marks
                    table fm
                    operator gte
                    value 75
                term twelfth_merit
                    column twelfth_percentage_marks
                    table fm
                    operator gte
                    value 75    
        
        evaluation parent_bocw && family_role && merit
    scheme mini_mata
        label Mini Mata Mahtari Jatan Yojna
        description """Female construction worker gets INR 20,000 during pregnancy"""
        criteria gender
            column gender
            table fm
            operator IN
            value ['female', 'other']
        criteria pregnancy
            column pregnancy_status
            table fm
            operator equals
            value true
        criteria bocw_card
            column has_bocw_card
            table fm
            operator equals
            value true
        criteria bocw_card_issue_date
            column bocw_card_issue_date
            table fm
            operator age_gte
            value 90
            granularity day
        
        evaluation gender && pregnancy && bocw_card && bocw_card_issue_date
    scheme naunihal
        label Naunihal Scholarship Scheme
        description """Children of construction workers get INR 500 - INR 5,000 per annum all the way from Class 1 to Postgraduate studies"""        
        criteria parent_bocw
            combine OR
                term mother_bocw         
                    column mother_bocw
                    table fm
                    operator equals
                    value True
                term father_bocw
                    column father_bocw
                    table fm
                    operator equals
                    value True
        criteria family_role
            column family_role
            table fm
            operator equals
            value child
        criteria in_educational_institute
            table fm
            column in_educational_institute
            operator equals
            value True
        
        evaluation parent_bocw && family_role && in_educational_institute
    scheme noni_shashaktikaran
        label Noni Sashaktikaran Scheme
        description """First 2 daughters in a family get INR 20,000 is directly transferred to bank account"""        
        criteria parent_bocw        
            combine OR
                term mother_bocw      
                    column mother_bocw
                    table fm
                    operator equals
                    value True
                term father_bocw
                    column father_bocw
                    table fm
                    operator equals
                    value True
        criteria family_rank
            column family_rank
            table fm
            operator IN
            value ['g1', 'g2']
        criteria age
            table fm
            column dob
            operator age_between
            value [18, 21]
            granularity year
        
        evaluation parent_bocw && family_rank && age
    scheme pm_ujjwala
        label Pradhan Mantri Ujjwala Yojana
        description """Free gas cylinders and stoves are provided to registered women laborers or registered construction laborers through the Food Department to the laborer's wife."""        
        criteria bocw
            combine OR
                term bocw_card
                    column has_bocw_card
                    table fm
                    operator equals
                    value True
                term father_bocw
                    column father_bocw
                    table fm
                    operator equals
                    value True
        criteria gender
            column gender
            table fm
            operator equals
            value female
        
        evaluation gender && bocw
    scheme sewing_machine
        label Chief Minister Sewing Machine Assistance Scheme
        description """Female construction worker receives sewing machine"""    
        criteria bocw_card
            column has_bocw_card
            table fm
            operator equals
            value true
        criteria age       
            column dob
            table fm
            operator age_between
            value [18, 50]
            granularity year
        criteria gender
            column gender
            table fm
            operator equals
            value female
        
        evaluation bocw_card && age && gender
    scheme shramik
        label Chief Minister Shramik Tool Assistance Scheme
        description """A one-time benefit of Rs. 1000 is provided to buy a tool kit (if tool kit is available then the kit is given)"""    
        criteria bocw_card
            column has_bocw_card
            table fm
            operator equals
            value true
        criteria age
            column dob
            table fm
            operator age_between
            value [18, 50]
            granularity year
        
        evaluation bocw_card && age
    scheme uow_cycle
        label Chief Minister Unorganized Workers Cycle Assistance Scheme
        description """One cycle (per beneficiary) will be payable."""    
        criteria uow_card
            column has_uow_card
            table fm
            operator equals
            value true
        criteria age
            column dob
            table fm
            operator age_between
            value [18, 40]
            granularity year
        criteria gender
            column gender
            table fm
            operator equals
            value female
        
        evaluation bocw_card && age && gender
            
				
`
	w, e := ParseSML(input, "")
	if e != nil {
		t.Errorf("Parse returned error: %s", e[0].Msg)
		return
	}
	b, e1 := json.Marshal(w)
	if e1 != nil {
		t.Errorf("json Marshal returned error: %v", e1)
		return
	}

	expectedJSON, _ := ioutil.ReadFile("testdata/smlCGRTEPlus.json")
	diff := DiffJSON(string(expectedJSON), string(b))
	if diff != "" {
		fmt.Printf("%v\n", string(b))
		t.Errorf("obtained json is not expected; diff=%s", diff)
	}
}
