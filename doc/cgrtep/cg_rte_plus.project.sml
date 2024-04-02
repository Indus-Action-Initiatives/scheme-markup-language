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
        column gender
            type string
            sql gender
        column family_role
            type string
            transformer lambda
            name findFamilyRole
        column relationship_with_respondent
            type string
            sql relationshipWithRespondent    
        column disadvantaged
            type bool
            sql disadvantaged
        column pregnancy_status
            type status
            transformer lambda
            name populatePregnancyStatus
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
            