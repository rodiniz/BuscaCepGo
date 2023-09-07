function addRows(data) {
    var table = document.getElementById('widget');
    if (!data) return;
    var tbody = table.createTBody();
    tbody.id = data.id;
    
    data.forEach(d => {
   
     var tr = tbody.insertRow();
     var td1 = tr.insertCell();
     var td2 = tr.insertCell();
     var td3 = tr.insertCell();
   
     var sz = document.createTextNode(d.Uf);
     td2.appendChild(td1);

     var sz = document.createTextNode(d.Logradouro);
     td2.appendChild(sz);
   
     var ph = document.createTextNode(d.Cep);
     td3.appendChild(ph);

    });
   }