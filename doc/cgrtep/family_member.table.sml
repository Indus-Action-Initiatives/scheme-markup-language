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

    # auxiliary columns start here
    auxiliary_column mother_bocw
        type bool
        sql SELECT mother.has_bocw_card FROM family_members as mother INNER JOIN families ON mother.family_id = families.id WHERE mother.family_role='mother' GROUP BY 1 LIMIT 1
    auxiliary_column father_bocw
        type bool
        sql SELECT father.has_bocw_card FROM family_members as father INNER JOIN families ON  father.family_id = fm.family_id WHERE father.family_role='father' GROUP BY 1 LIMIT 1    
