$(document).ready(function(){
	var nocm = $("#nocm").val()
	
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
		var ats = $("input[name='ats']:checked").val();
		var iki = $("input[name='iki']:checked").val();
		var shift = $("input[name='shift']:checked").val();
		var bagian = $("input[name='bagian']:checked").val();
		var baru = $("#baru").val();
		var dok = $("#email").val();

		switch (true) {
			case nocm === "":
			break;
			case namapts === "":
			break;
			case diag === "":
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Diagnosis belum diisi" +
		                        "</div>");
			// alert("Diagnosis harus diisi");
			break;
			case ats == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"ATS belum diisi" +
		                        "</div>");
			// alert("ATS harus diisi");
			break;
			case iki == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Golongan IKI belum diisi" +
		                        "</div>");
			// alert("Golongan IKI harus diisi");
			break;
			case shift == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Shift belum diisi" +
		                        "</div>");
			// alert("Shift harus diisi");
			break;
			case bagian == null:
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Bagian belum diisi" +
		                        "</div>");
			// alert("Bagian harus diisi");
			break;
			default:
			if (namapts !== "" && diag !== "" && ats !== undefined && iki !== undefined && shift !== undefined && bagian !== undefined) {
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
				var js = JSON.parse(data)
				if (js.token != "OK"){
					alert(js.script);
				}else{
  				    $("tbody").prepend(js.script);
					  alert("Data berhasil diinput");
					  $("#nocm").val('');
					  $("#datapasien").html('');
				}


			});
			}else{
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Ada kolom belum diisi" +
		                        "</div>");
			}
		}
		// alert("Form belum lengkap!")		
	});

	$("#navbar").on("click", "#editbut", function(event){
		event.preventDefault();
		token = localStorage.getItem("token")
		var link = $(this).offsetParent().children().first().html();
		// var baris = $(this).parent("ul").index();
		// console.log("Ini baris ke: " + baris);
		$.post("editentri",{
			token: token,
			link: link
		},
		function(data){
			var js = JSON.parse(data)
			$(js.script).modal();
			$("body").on('shown.bs.modal', function(){
			var ats = 'input[name="modats"][value=' + js.content.ats + ']';
			var iki = 'input[name="modiki"][value=' + js.content.iki + ']';
			var shift = 'input[name="modshift"][value=' + js.content.shift + ']';
			var bagian = 'input[name="modbagian"][value=' + js.content.bagian + ']';
			
			$('input[name="entri"]').val(js.content.link);
			$('input[name="namapasien"]').val(js.content.namapts);
			$('input[name="diagnosis"]').val(js.content.diag);
			
			if (js.content.bagian == ""){
 		    	$(ats).prop('checked', true);
   				$(iki).prop('checked', true);
   				$(shift).prop('checked', true);	
			
			}else{
 		    	$(ats).prop('checked', true);
   				$(iki).prop('checked', true);
   				$(shift).prop('checked', true);
   				$(bagian).prop('checked', true);

				};
			});
		})
	});

	$("#modedit").on("click", "#confirmedit", function(){

		var link = $("#modentri").val();
		var namapts = $("#modnamapts").val();
		var diag = $("#moddiag").val();
		var ats = $("input[name='modats']:checked").val();
		var bagian = $("input[name='modbagian']:checked").val();
		var iki = $("input[name='modiki']:checked").val();
		var shift = $("input[name='modshift']:checked").val();
		
		switch (true) {
			case namapts === "":
			$("#alertmodal").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Nama pasien belum diisi" +
		                        "</div>");
			break;
			case diag === "":
			$("#alertmodal").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Diagnosis belum diisi" +
		                        "</div>");
			// alert("Diagnosis harus diisi");
			break;
			case ats == null :
			$("#alertmodal").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"ATS belum diisi" +
		                        "</div>");
			// alert("ATS harus diisi");
			break;
			case iki == null :
			$("#alertmodal").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Golongan IKI belum diisi" +
		                        "</div>");
			// alert("Golongan IKI harus diisi");
			break;
			case shift == null :
			$("#alertmodal").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Shift belum diisi" +
		                        "</div>");
			// alert("Shift harus diisi");
			break;
			case bagian == null:
			$("#alertmodal").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Bagian belum diisi" +
		                        "</div>");
			// alert("Bagian harus diisi");
			break;
			default:
			if (namapts !== "" && diag !== "" && ats !== undefined && iki !== undefined && shift !== undefined && bagian !== undefined) {
			$.post("confedit",{
				token: localStorage.getItem("token"),
				nocm: nocm,
				namapts: namapts,
				diag: diag,
				ats: ats,
				iki: iki,
				shift: shift,
				bagian: bagian,
				link: link
			},
			function(data){
				var js = JSON.parse(data)
				if (js.token != "OK"){
					alert(js.script);
				}else{
  				    $(this).parent("tr")(js.script);
					  alert("Data berhasil diubah");
					//   $("#nocm").val('');
					//   $("#datapasien").html('');
				}


			});
			}else{
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Ada kolom belum diisi" +
		                        "</div>");
			}
		}
		


	});

})