import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'nav/customer_nav.dart';
import 'homecustomer.dart';

class BuyerProfile extends StatefulWidget {
  final String userId;
  const BuyerProfile({Key? key, required this.userId}) : super(key: key);

  @override
  _BuyerProfileState createState() => _BuyerProfileState();
}

class _BuyerProfileState extends State<BuyerProfile> {
  String id = '';
  String firstName = '';
  String lastName = '';
  String email = '';
  String mobile = '';

  @override
  void initState() {
    super.initState();
    fetchUserData();
  }

  Future<void> fetchUserData() async {
  final url = Uri.parse('http://localhost:3000/v1/customers/?id=${widget.userId}');
  try {
    final response = await http.get(url);
    if (response.statusCode == 200) {
      final data = json.decode(response.body);

      setState(() {
        id = data['cust_id']?.toString() ?? 'N/A';
        firstName = data['cust_fname'] ?? 'N/A';
        lastName = data['cust_lname'] ?? 'N/A';
        email = data['cust_email'] ?? 'N/A';
      });
    } else {
      print('Failed to load user data: ${response.statusCode}');
    }
    } catch (e) {
      print('Error fetching user data: $e');
    }
  }



  @override
  Widget build(BuildContext context) {
    return Scaffold(
     appBar: AppBar(
        toolbarHeight: 80, 
        title: Row(
          children: [
            Icon(Icons.location_on, color: Colors.red, size: 35,),
            SizedBox(width: 8),
            Text(
              'My Account',
              style: TextStyle(color: Colors.red, fontSize: 24, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        automaticallyImplyLeading: false,
        actions: [
          Container(
            margin: EdgeInsets.symmetric(horizontal: 25.0),
            child: IconButton(
              icon: Icon(Icons.chevron_left, color: Colors.red, size: 35),
              onPressed: () {
                Navigator.pushReplacement(
                  context,
                  MaterialPageRoute(builder: (context) => HomeScreen(userId: widget.userId)),
                );
              },
            ),
          ),
        ],
      ),
    body: id.isEmpty
    ? const Center(child: CircularProgressIndicator())
    : Padding(
  padding: const EdgeInsets.all(16.0),
  child: Column(
    crossAxisAlignment: CrossAxisAlignment.start,
    children: [
      Container(
        width: double.infinity,
        padding: const EdgeInsets.only(bottom: 10),
        decoration: const BoxDecoration(
          border: Border(
            bottom: BorderSide(color: Colors.grey, width: 1),
          ),
        ),
        child: Text(
          'ID: $id',
          style: const TextStyle(color: Colors.black, fontSize: 20),
        ),
      ),
      const SizedBox(height: 15),
      const Text(
        'First Name',
        style: TextStyle(color: Colors.black, fontSize: 20),
      ),
      const SizedBox(height: 5),
      Container(
        width: double.infinity,
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          border: Border.all(color: Colors.grey),
          borderRadius: BorderRadius.circular(5),
        ),
        child: Text(
          firstName,
          style: const TextStyle(color: Colors.black, fontSize: 16),
        ),
      ),
      const SizedBox(height: 15),
      const Text(
        'Last Name',
        style: TextStyle(color: Colors.black, fontSize: 20),
      ),
      const SizedBox(height: 5),
      Container(
        width: double.infinity,
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          border: Border.all(color: Colors.grey),
          borderRadius: BorderRadius.circular(5),
        ),
        child: Text(
          lastName,
          style: const TextStyle(color: Colors.black, fontSize: 16),
        ),
      ),
      const SizedBox(height: 15),
      const Text(
        'Email',
        style: TextStyle(color: Colors.black, fontSize: 20),
      ),
      const SizedBox(height: 5),
      Container(
        width: double.infinity,
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          border: Border.all(color: Colors.grey),
          borderRadius: BorderRadius.circular(5),
        ),
        child: Text(
          email,
          style: const TextStyle(color: Colors.black, fontSize: 16),
        ),
      ),
    ],
  ),
),

      bottomNavigationBar: BottomNavBar(userId: widget.userId),
    );
  }
}
