import 'package:bazaroo_client/address.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'nav/customer_nav.dart';

class AddressForm extends StatefulWidget {
  final String userId;
  final int addrId;
  AddressForm({required this.userId, required this.addrId});

  @override
  State<AddressForm> createState() => _NewAddressState();
}

class _NewAddressState extends State<AddressForm> {
  final TextEditingController addressLine1Controller = TextEditingController();
  final TextEditingController addressLine2Controller = TextEditingController();
  final TextEditingController cityController = TextEditingController();
  final TextEditingController stateController = TextEditingController();
  final TextEditingController postalCodeController = TextEditingController();
  final TextEditingController countryController = TextEditingController();
  final TextEditingController phoneController = TextEditingController();
  bool isLoading = false;
  String? errorMessage;

  Future<void> updateAddr(int addrId) async {
    final url = Uri.parse("http://localhost:3000/v1/customer/addr/?id=${widget.userId}");
    http.put(url, headers: {"Content-Type": "application/json"},
      body: jsonEncode({
        "addr_id": addrId,
      }),
    );
  }
  Future<void> registerBusiness() async {
    setState(() {
      isLoading = true;
      errorMessage = null;
    });

    final url = Uri.parse("http://localhost:3000/v1/addr");
    final response = await http.post(
      url,
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({
        "addr_line1": addressLine1Controller.text,
        "addr_line2": addressLine2Controller.text,
        "city": cityController.text,
        "state": stateController.text,
        "postal_code": postalCodeController.text,
        "country": countryController.text,
      }),
    );

    if (response.statusCode == 201) {
      final resMess = jsonDecode(response.body);
      var message = resMess['addr_id'];
      print(message);
       Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => AddressScreen(userId: widget.userId)),
      );
    } else {
      setState(() {
        errorMessage = "Failed to register address. Try again.";
      });
    }

    setState(() {
      isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        toolbarHeight: 80,
        backgroundColor: Colors.red,
        title: Text(
          'Address Form',
          style: TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
        ),
        automaticallyImplyLeading: false,
        actions: [
          Container(
            margin: EdgeInsets.symmetric(horizontal: 25.0),
            child: IconButton(
              icon: Icon(Icons.chevron_right, color: Colors.white, size: 35),
              onPressed: () {
                Navigator.pop(context);
              },
            ),
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            SizedBox(height: 30),
            if (errorMessage != null)
              Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: Text(
                  errorMessage!,
                  style: TextStyle(color: Colors.red, fontSize: 16),
                ),
              ),
            TextFormField(
              controller: addressLine1Controller,
              decoration: InputDecoration(
                labelText: 'Address Line 1',
                prefixIcon: const Icon(Icons.location_on),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
            ),
            const SizedBox(height: 10),
            TextFormField(
              controller: addressLine2Controller,
              decoration: InputDecoration(
                labelText: 'Address Line 2 (Optional)',
                prefixIcon: const Icon(Icons.location_on_outlined),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
            ),
            const SizedBox(height: 10),
            TextFormField(
              controller: cityController,
              decoration: InputDecoration(
                labelText: 'City',
                prefixIcon: const Icon(Icons.location_city),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
            ),
            const SizedBox(height: 10),
            TextFormField(
              controller: stateController,
              decoration: InputDecoration(
                labelText: 'State/Province',
                prefixIcon: const Icon(Icons.map),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
            ),
            const SizedBox(height: 10),
            TextFormField(
              controller: postalCodeController,
              decoration: InputDecoration(
                labelText: 'Postal Code',
                prefixIcon: const Icon(Icons.local_post_office),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
            ),
            const SizedBox(height: 10),
            TextFormField(
              controller: countryController,
              decoration: InputDecoration(
                labelText: '3 Country Code',
                prefixIcon: const Icon(Icons.flag),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
            ),
            const SizedBox(height: 10),
            SizedBox(height: 30),
            ElevatedButton(
              onPressed: isLoading ? null : registerBusiness,
              style: ElevatedButton.styleFrom(
                minimumSize: Size(double.infinity, 50),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
                backgroundColor: Colors.red,
              ),
              child: isLoading
                  ? CircularProgressIndicator(color: Colors.white)
                  : Text(
                      'Register',
                      style: TextStyle(fontSize: 18, color: Colors.white),
                    ),
            ),
          ],
        ),
      ),
    );
  }
}
