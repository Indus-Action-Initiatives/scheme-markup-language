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
        