$(document).ready(function(){
	var nocm = $("#nocm").val()
	
	// $("#nocm").focus(function(){
	// 	var value = $("#nocm").val();
	// 	if (value == ""){
	// 		$("#datapasien").html("Masukkan No. CM");
	// 		//nocm = "";
	// 	} else {
	// 		$("#datapasien").html("No. CM tidak lengkap");
	// 		//nocm = "";
	// 	} 
	// });
	
	$("#navbar").on("keyup", "#nocm", function(){
		
		var value = $("#nocm").val();
		
		
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
		} else if (value.length < 8){
			$("#datapasien").html("No. CM tidak lengkap");
		} else {
			$("#nocm").prop("disabled", true);
			nocm = value;
            token = localStorage.getItem("token");

			$.post("getcm", {
				nocm: nocm,
				token: token
			},
			function(data){
				var js = JSON.parse(data)

				if (js.token == "OK"){
					$("#datapasien").html(js.script);
					$("#nocm").prop("disabled", false);
				} else {
					alert("Terjadi kesalahan: " + js.token)
					$("#nocm").prop("disabled", false);
				}			    

			});
        }

	});
	
	$("#navbar").on("click", "#btnsub", function(event){
		event.preventDefault();
		
		var nocm = $("#nocm").val();
		var namapts = $("#namapts").val();
		var diag = $("#diag").val();
		var ats = $("input[type='radio'][name='ats']:checked").val();
		var iki = $("input[type='radio'][name='iki']:checked").val();
		var shift = $("input[type='radio'][name='shift']:checked").val();
		var bagian = $("input[type='radio'][name='bagian']:checked").val();
		var baru = $("#baru").val();
		var dok = $("#email").val();
		console.log("Dokter adalah: " + dok)
		if (nocm !== "" && namapts !== "" && diag !== "" && ats !== "" && iki !== "" && shift !== "" && bagian !== "" ){
			$.post("inputdata",{
				token: localStorage.getItem("token"),
				nocm: nocm,
				namapts: namapts,
				diag: diag,
				ats: ats,
				iki: iki,
				shift: shift,
				bagian: bagian,
				dok: dok,
				baru: baru
			},
			function(data){
				$("#nocm").empty();
				var js = JSON.parse(data)
				if (js.token != "OK"){
					alert(js.script)
				}else{
  				    $("tbody").prepend(js.script)
				}
			})
		}else{
			alert("Data Belum Lengkap");
		}
		
	});

	$("#navbar").on("click", "#editbut", function(event){
		event.preventDefault();
		token = localStorage.getItem("token")
		var link = $(this).offsetParent().children().first().html();
		console.log("link adalah: " + link)
		$.post("editentri",{
			token: token,
			link: link
		},
		function(data){
			var js = JSON.parse(data)
			console.log(js.token)
			$(js.script).appendTo('body').modal()
		})

	})

})