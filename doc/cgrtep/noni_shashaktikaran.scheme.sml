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