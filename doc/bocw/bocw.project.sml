project bocw
    dataset main
        table fm
    table fm
        label fm
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