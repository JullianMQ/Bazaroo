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
  Map<String, dynamic>? products;
  bool isLoading = true; // Track loading state

  Future<void> delAddresses() async {
    final putUrl = Uri.parse('http://localhost:3000/v1/customers/addr/?id=${widget.userId}');
    final delUrl = Uri.parse('http://localhost:3000/v1/addr/?id=${widget.userId}');

      final response = await http.put(putUrl);
      Map<String, dynamic> jsonData = json.decode(response.body);
      setState(() {
        products = jsonData;
        products?['addr_id']=0;
        http.delete(delUrl);
      });
  }

  Future<void> fetchAddresses() async {
    final url = Uri.parse('http://localhost:3000/v1/customers/?id=${widget.userId}');
    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {
        Map<String, dynamic> jsonData = json.decode(response.body);
        setState(() {
          products = jsonData;
          isLoading = false;
        });
      } else {
        throw Exception('Failed to load products: ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching products: $e');
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  void initState() {
    super.initState();
    fetchAddresses();
  }

  @override
  Widget build(BuildContext context) {
    int addrId = products?['addr_id']?['Int64'] ?? 0; // Safely get addrId

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
                            '${products?['cust_fname'] ?? ''} ${products?['cust_lname'] ?? ''}',
                            style: TextStyle(fontWeight: FontWeight.bold, color: Colors.black87),
                          ),
                        ),
                        Text(
                          '${products?['phone_num'] ?? ''}',
                          style: TextStyle(color: Colors.grey),
                        ),
                      ],
                    ),
                    Text(
                      '## St. ********************',
                      style: TextStyle(color: Colors.grey),
                    ),
                    Text(
                      'Pampanga, Angeles, *****',
                      style: TextStyle(color: Colors.grey),
                    ),
                    SizedBox(height: 8),

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
                          onTap: (delAddresses),
                          child: Text(
                            'Delete',
                            style: TextStyle(color: Colors.red, decoration: TextDecoration.underline),
                          ),
                        ),
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
