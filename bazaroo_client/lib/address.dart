import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'AddressForm.dart';
import 'homecustomer.dart';

class AddressScreen extends StatefulWidget {
  final String userId;
  AddressScreen({required this.userId});

  @override
  State<AddressScreen> createState() => _AddressScreenState();
}

class _AddressScreenState extends State<AddressScreen> {
  Map<String, dynamic>? customer;
  String address ="";
  bool isLoading = true;

  Future<void> updId() async {
    final putUrl = Uri.parse('https://bazaroo.onrender.com/v1/customers/addr/?id=${widget.userId}');
      http.put(putUrl, headers: {"Content-Type": "application/json"},
        body: jsonEncode({"addr_id":0}),
      );
  }

  Future<void> delAddresses(String id) async {
    try {
      updId();
      final delUrl = Uri.parse('https://bazaroo.onrender.com/v1/addr/?id=$id');
      final response = await http.delete(delUrl);

      if (response.statusCode == 201) {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (context) => AddressScreen(userId: widget.userId)),
        );
      } else {
        throw Exception('Failed to delete address: ${response.body}');
      }
    } catch (e) {
      print(e); 
    }
  }


  Future<void> fetchCustomers() async {
    final url = Uri.parse('https://bazaroo.onrender.com/v1/customers/?id=${widget.userId}');
    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {

        Map<String, dynamic> jsonData = json.decode(response.body);
        setState(() {
          customer = jsonData;
          int addrId = customer?['addr_id']["Int64"] ?? 0;
          fetchAddresses(addrId);
          isLoading = false;
        });
      } else {
        throw Exception('Failed to load customer: ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching customer: $e');
      setState(() {
        isLoading = false;
      });
    }
  }

  Future<void> fetchAddresses(int addrId) async {
    final url = Uri.parse('https://bazaroo.onrender.com/v1/addr/?id=$addrId');
    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {
        Map<String, dynamic> jsonData = json.decode(response.body);
        setState(() {
          Map<String, dynamic> content = jsonData;
          address = '${content['addr_line1']} ${content['addr_line2']['String']} ${content['city']} ${content['state']} ${content['country']} ${content['postal_code']}' ;    
          isLoading = false;
        });
      } else {
        throw Exception('Failed to load address: ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching address: $e');
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  void initState() {
    super.initState();
    fetchCustomers();
  }

  @override
  Widget build(BuildContext context) {
    int addrId = customer?['addr_id']?['Int64'] ?? 0; 

    return Scaffold(
      appBar: AppBar(
        toolbarHeight: 80,
        title: Row(
          children: [
            Icon(Icons.location_on, color: Colors.red, size: 35),
            SizedBox(width: 8),
            Text(
              'Address',
              style: TextStyle(color: Colors.red, fontSize: 24, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        automaticallyImplyLeading: false,
        actions: [
          Container(
            margin: EdgeInsets.symmetric(horizontal: 25.0),
            child: IconButton(
              icon: Icon(Icons.chevron_right, color: Colors.red, size: 35),
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
      body: isLoading
          ? Center(child: CircularProgressIndicator()) 
          : Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (addrId != 0) ...[
                    Center(
                      child: ElevatedButton(
                        onPressed: null, 
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.grey[400],
                          foregroundColor: Colors.white,
                          padding: EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(10),
                          ),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.add, size: 20, color: Colors.white),
                            SizedBox(width: 8),
                            Text('Add Address'),
                          ],
                        ),
                      ),
                    ),
                    SizedBox(height: 16),
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            '${customer?['cust_fname'] ?? ''} ${customer?['cust_lname'] ?? ''}',
                            style: TextStyle(fontWeight: FontWeight.bold, color: Colors.black),
                          ),
                        ),
                        Text(
                          '${customer?['phone_num'] ?? ''}',
                          style: TextStyle(color: Colors.grey),
                        ),
                      ],
                    ),
                    Text(
                      address,
                      style: TextStyle(color: Colors.grey),
                    ),
                    Row(
                      children: [
                        InkWell(
                          onTap: () { 
                           Navigator.push(
                            context,
                            MaterialPageRoute(
                              builder: (context) => AddressForm(userId: widget.userId, addrId: addrId),
                            )
                           );
                          },
                          child: Text(
                            'Edit',
                            style: TextStyle(color: Colors.red, decoration: TextDecoration.underline),
                          ),
                        ),
                        SizedBox(width: 16),
                        InkWell(
                          onTap: () => delAddresses((customer?['addr_id']['Int64'] ?? '').toString()),
                          child: Text(
                            'Delete',
                            style: TextStyle(color: Colors.red, decoration: TextDecoration.underline),
                          ),
                        )
                      ],
                    ),
                    Divider(thickness: 1, color: Colors.grey[300]),
                  ] else ...[
                    Center(
                      child: ElevatedButton(
                        onPressed: () {
                          Navigator.push(
                            context,
                            MaterialPageRoute(
                              builder: (context) => AddressForm(userId: widget.userId, addrId: addrId),
                            ),
                          );
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.red,
                          foregroundColor: Colors.white,
                          padding: EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(10),
                          ),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.add, size: 20, color: Colors.white),
                            SizedBox(width: 8),
                            Text('Add Address'),
                          ],
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),
    );
  }
}
