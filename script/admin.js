$(document).ready(function(){
    $("#navbar").on("click", "#admbtnsub", function(e){
        e.preventDefault();
        namdok = $("input#namadktr").val();
        email = $("input#emaildktr").val();
        peran = $("#perandktr option:checked").val();
		if (peran == "residen"){
			var residen = $("#adm-select-residen option:checked").val();
			peran = peran + "-" + residen
		} else {
			peran = peran
		}
        $.post("tambahdokter", {
            token: localStorage.getItem("token"),
            nama: namdok,
            email: email,
            peran: peran
        }, function(data){
            js = JSON.parse(data)
            $("input").val("");
            $("tbody#staftab").append(js.script)
            $("#mymodal").html(js.modal);
            $("tr").find(".nourut").each(function(index, elem){
                num = index + 1;
                $(elem).html(num)
            });
            $("#mymodal").modal();
        })
    })

    $("#navbar").on("click", "a.deldok", function(e){
        e.preventDefault();
        link = $(this).attr("href")
        indexrow = $(this).parents("tr").index() +1 ;
        console.log(indexrow)
        $.post("hapusdokter", {
            token: localStorage.getItem("token"),
            link: link
        }, function(data){
            js = JSON.parse(data)
            if (js.token == "OK"){
                // indexrow = $(this).closest("tr").index();
                urutan = "tr:eq(" + indexrow + ")"
                $(urutan).remove()
                $("tr").find(".nourut").each(function(index, elem){
                    num = index + 1;
                    $(elem).html(num)
                });
                $("#mymodal").html(js.modal);
                $("#mymodal").modal();
                
            }else{
                $("#mymodal").html(js.modal)
                $("#mymodal").modal();
            }
        })
    })
	$("#navbar").on("change", "#perandktr", (function(){
		// console.log("triggered")
		var peran = $("#perandktr").val()
		// console.log("Peran adalah: " + peran)
		if (peran == "residen"){
			$("#div-adm-select-residen").show()
		}else{
			$("#div-adm-select-residen").hide()
		}
	}))
})