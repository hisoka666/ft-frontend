// $(document).ready(function(){
    // $("#jmlperhari").html("tolong isi ini!")
    function chartOn(){
        console.log("This function works")
        google.charts.load('current', {packages:['corechart']});
        google.charts.setOnLoadCallback(barHari);
        var pagi = []
        $(".shift-pagi").each(function(){
            console.log($(this).html())
            pagi.push($(this).html())
        })
        var sore = []
        $(".shift-sore").each(function(){
            console.log($(this).html())
            sore.push($(this).html())
        })
        var malam = []
        $(".shift-malam").each(function(){
            console.log($(this).html())
            malam.push($(this).html())
        })
        var total = []
        $(".shift-total").each(function(){
            console.log($(this).html())
            total.push($(this).html())
        })

        function barHari(){
            var data = new google.visualization.DataTable();
            data.addColumn('number', 'Tanggal');
            data.addColumn('number', 'Pagi');
            data.addColumn('number', 'Sore');
            data.addColumn('number', 'Malam');
            data.addColumn('number', 'Total');
    
            
            var isi = []
            for(i=0;i<31;i++){
                data.addRows([
                    [i+1,pagi[i], sore[i], malam[i], total[i]]
                ])
            }

            var options = {
                title: "Jumlah Kunjungan IGD Perhari",
                width: 1000,
                hAxis: {
                    title: "Tanggal Jaga",
                },
                vAxis: {
                    title: "Jumlah Pasien"
                }
            };

            var chart = new google.visualization.ColumnChart($("#jmlperhari"));
            chart.draw(data, options);
        }
    }
    // $("#navbar").on("click", "#jmlperhari", chartOn())
    // chartOn();
    // jmlPerHari();
    // function jmlPerHari(){
        
    // }
    // })