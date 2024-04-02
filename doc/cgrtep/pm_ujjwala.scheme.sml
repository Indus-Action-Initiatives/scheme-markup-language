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