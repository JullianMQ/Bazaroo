import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'homeseller.dart';
import 'loginseller.dart';

class RegisterBusiness extends StatefulWidget {
  final String userId;
  const RegisterBusiness({Key? key, required this.userId}) : super(key: key);

  @override
  _RegisterBusinessState createState() => _RegisterBusinessState();
}

class _RegisterBusinessState extends State<RegisterBusiness> {
  final TextEditingController addressLine1Controller = TextEditingController();
  final TextEditingController addressLine2Controller = TextEditingController();
  final TextEditingController cityController = TextEditingController();
  final TextEditingController stateController = TextEditingController();
  final TextEditingController postalCodeController = TextEditingController();
  final TextEditingController countryController = TextEditingController();
  final TextEditingController phoneController = TextEditingController();
  String? errorMessage;
  int? officeId;
  int? addrId;

  bool isLoading = false;

  @override
  void dispose() {
    addressLine1Controller.dispose();
    addressLine2Controller.dispose();
    cityController.dispose();
    stateController.dispose();
    postalCodeController.dispose();
    countryController.dispose();
    phoneController.dispose();
    super.dispose();
  }

  Future<void> addOfficeId(int id) async {
    final url = Uri.parse("http://localhost:3000/v1/emps/?id=${widget.userId}");

    try {
      final response = await http.put(
        url,
        headers: {"Content-Type": "application/json"},
        body: jsonEncode({"office_id": id}),
      );

      if (response.statusCode != 201) {
        throw Exception('Failed to add office id: $response.statusCode');
      }

      setState(() {
        officeId = id;
        isLoading = true;
      });
    } catch (e) {
      throw Exception('Error adding office id: $e');
    }
  }

  Future<void> addAddrId(int id) async {
    final url = Uri.parse("http://localhost:3000/v1/offices");

    try {
      final res = await http.post(
        url,
        headers: {"Content-Type": "application/json"},
        body: jsonEncode({
          "phone_num": phoneController.text.trim(),
          "addr_id": id}),
      );

      if (res.statusCode != 201) {
        throw Exception('Failed to add office id: $res.statusCode');
      }

      final responseData = jsonDecode(res.body);

      setState(() {
        officeId = int.parse(responseData["message"]);
        isLoading = true;
      });
    } catch (e) {
      throw Exception('Error adding office id: $e');
    }
  }

  Future<void> registerBusiness() async {
    if (addressLine1Controller.text.isEmpty ||
        cityController.text.isEmpty ||
        stateController.text.isEmpty ||
        postalCodeController.text.isEmpty ||
        countryController.text.isEmpty ||
        phoneController.text.isEmpty) {
      setState(() {
        errorMessage = 'Please fill in all required fields';
      });
      return;
    }

    setState(() {
      isLoading = true;
      errorMessage = null;
    });

    final url = Uri.parse("http://localhost:3000/v1/addr");

    try {
      final response = await http.post(
        url,
        headers: {"Content-Type": "application/json"},
        body: jsonEncode({
          "addr_line1": addressLine1Controller.text.trim(),
          "addr_line2": addressLine2Controller.text.trim(),
          "city": cityController.text.trim(),
          "state": stateController.text.trim(),
          "postal_code": postalCodeController.text.trim(),
          "country": countryController.text.trim(),
        }),
      );

      if (response.statusCode == 201) {
        final responseData = jsonDecode(response.body);

        try {
          addrId = int.parse(responseData["message"]);
            await addAddrId(addrId!);
            await addOfficeId(officeId!);
        } catch (e) {
          setState(() {
            errorMessage = e.toString();
          });
          return;
        }

        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
              builder: (context) => SellerHome(userId: widget.userId)),
        );
      } else {
        final responseData = jsonDecode(response.body);
        setState(() {
          errorMessage = responseData["error"] ??
              "Unexpected Error: ${response.statusCode}";
        });
      }
    } catch (e) {
      setState(() {
        errorMessage = 'Error: $e';
      });
    } finally {
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        automaticallyImplyLeading: false,
        backgroundColor: Colors.white,
        toolbarHeight: 80,
        title: Row(
          children: [
            IconButton(
              icon: Icon(Icons.chevron_right, color: Colors.red, size: 35),
              onPressed: () {
                Navigator.pushReplacement(
                  context,
                  MaterialPageRoute(builder: (context) => LoginSeller()),
                );
              },
            ),
            Text(
              'Register Business',
              style: TextStyle(color: Colors.red, fontSize: 40),
            ),
          ],
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Image.asset(
              'assets/Bazaroo-full-logo-red.png',
              height: 300,
            ),
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
            TextFormField(
              controller: phoneController,
              decoration: InputDecoration(
                labelText: '09xxxxxxxxxx',
                prefixIcon: const Icon(Icons.phone),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                ),
              ),
            ),
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
