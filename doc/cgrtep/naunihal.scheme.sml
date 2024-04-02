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
        