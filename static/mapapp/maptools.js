var selectedGeom
var currentGeometry
var currentMarker
var osmUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
    osmAttrib = '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    osm = L.tileLayer(osmUrl, { maxZoom: 18, attribution: osmAttrib }),
    map = new L.Map('map', { center: new L.LatLng(37.1, -90.4), zoom: 5 }),
    drawnItems = L.featureGroup().addTo(map);
L.control.layers({
    'osm': osm.addTo(map),
    "google": L.tileLayer('http://www.google.cn/maps/vt?lyrs=s@189&gl=cn&x={x}&y={y}&z={z}', {
        attribution: 'google'
    })
}, { 'drawlayer': drawnItems }, { position: 'topleft', collapsed: false }).addTo(map);
map.addControl(new L.Control.Draw({
    edit: {
        featureGroup: drawnItems,
        poly: {
            allowIntersection: false
        }
    },
    draw: {
        polygon: {
            allowIntersection: false,
            showArea: true,
            shapeOptions: {
                color: 'purple'
            }
        },
        circle: {
            showArea: true,
            shapeOptions: {
                color: 'red'
            }
        }
    }
}));
var findLayer = L.geoJSON(false, { onEachFeature: bindPopupOnEachFeature}).addTo(map);

function zoomTo() {
    var latlng = document.getElementById("latlngform").value;
    splitted = latlng.split(',')
    var lat = Number(splitted[0])
    var lng = Number(splitted[1])
    map.flyTo(new L.LatLng(lat, lng), 17);
    if (currentMarker != undefined) {
        map.removeLayer(currentMarker);
    }
    currentMarker = L.marker([lat,lng]).addTo(map);
}

function find() {
    var store_id = document.getElementById("store_id_input").value;
    var metro_id = document.getElementById("metro_id_input").value;
    var zone_id = document.getElementById("zone_id_input").value;
    var city = document.getElementById("city_input").value;
    var state = document.getElementById("state_input").value;
    var reqBody = JSON.stringify({"store_id": newParseInt(store_id), "metro_id": newParseInt(metro_id), "zone_id": newParseInt(zone_id), "city": city, "state": state})
    xhr = new XMLHttpRequest();
    var url = "/poly/find";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var json = JSON.parse(xhr.responseText);
            findLayer.clearLayers()
            findLayer.addData(json)
            map.flyTo(new L.LatLng(37.1, -90.4), 5)
        }
    }
    xhr.send(reqBody)
}

function findByID() {
    var id = document.getElementById("id_input").value;
    findByIDhelper(id)
}

function next() {
    if (currentGeometry) {
        console.log(currentGeometry)
        var id = currentGeometry.properties.ID
        findByIDhelper(id + 1)
        map.flyTo(new L.LatLng(currentGeometry.properties.Latitude, currentGeometry.properties.Longitude), 16);
    } else {
        findByIDhelper(1)
    }
}

function prev() {
    if (currentGeometry && currentGeometry.properties.ID != 1) {
        console.log(currentGeometry)
        var id = currentGeometry.properties.ID
        findByIDhelper(id - 1)
        map.flyTo(new L.LatLng(currentGeometry.properties.Latitude, currentGeometry.properties.Longitude), 16);
    } else {
        findByIDhelper(1)
    }
}

function findByIDhelper(id) {
    xhr = new XMLHttpRequest();
    var url = "/poly/find/" + id;
    xhr.open("GET", url, true);
    xhr.setRequestHeader("Content-type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var json = JSON.parse(xhr.responseText);
            findLayer.clearLayers()
            findLayer.addData(json)
            currentGeometry = json[0]
            map.flyTo(new L.LatLng(currentGeometry.properties.Latitude, currentGeometry.properties.Longitude), 16)
            renderPolygonIfAvailable(currentGeometry)
        }
    }
    xhr.send()
}

function renderPolygonIfAvailable(currentGeometry) {
    if (currentGeometry.properties.Polygon) {
        findLayer.addData(asJSONfeature(JSON.parse(currentGeometry.properties.Polygon), currentGeometry.properties))
    }
}

function asJSONfeature(geometry, properties) {
    return {"type": "Feature", "properties": properties, "geometry": geometry}
}
function newParseInt(stringNumber) {
    if (stringNumber == "") {
        return 0
    } else {
        return parseInt(stringNumber, 10)
    }
}

function bindPopupOnEachFeature(feature, layer) {
    // does this feature have a property named popupContent?
    if (feature.properties) {
        layer.bindTooltip(popupFromProperties(feature.properties));
        layer.on('click', function(e) {map.flyTo(new L.LatLng(feature.properties.Latitude, feature.properties.Longitude), 16);});
        currentGeometry = feature
    } else {
        layer.bindTooltip("No properties available")
    }
}

function popupFromProperties(properties) {
    idDiv = "<div><b>" + properties.ID + "</b></div>"
    nameDiv = "<div>" + properties.Name + "</div>"
    addrDiv = "<div>" + properties.Street1 + ", " + properties.City + ", " + properties.State + " " + properties.Zip + "</div>"
    return idDiv + nameDiv + addrDiv
}

function submitPolygon() {
    var id = Number(document.getElementById('ide').value)
    var geomObj = selectedGeom
    var reqBody = JSON.stringify({"id": id, "polygon": geomObj})
    console.log(reqBody)
    xhr = new XMLHttpRequest();
    var url = "/insert/poly";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var json = JSON.parse(xhr.responseText);
        }
    }
    xhr.send(reqBody);
    map.closePopup();
    return false
}

map.on(L.Draw.Event.CREATED, function (event) {
    var layer = event.layer;
    var geoJSON = event.layer.toGeoJSON();
    selectedGeom = geoJSON.geometry;
    var tempMarker = drawnItems.addLayer(layer);
    var ide = currentGeometry ? currentGeometry.properties.ID : "";
    var name = currentGeometry ? currentGeometry.properties.Name : "";
    var popupForm = `<form enctype="application/json">
                              <fieldset>
                                <legend>Store Polygon:</legend>
                                ID:<br>
                                <input type="number" name="id" value="` + ide + `" id="ide"><br>
                                <div>` + name + `<div>
                                <button type="button" onclick="submitPolygon(); return false;">Submit</button>
                              </fieldset>
                            </form>`
    var popupContent = popupForm
    layer.bindPopup(popupContent,{
        keepInView: false,
        closeButton: true
    }).openPopup();
    layer.on("click", function (e) {
        selectedGeom = layer.toGeoJSON().geometry
    })

});

