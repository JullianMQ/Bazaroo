import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'nav/customer_nav.dart';

class Cart extends StatefulWidget {
  final String userId;
  const Cart({Key? key, required this.userId}) : super(key: key);

  @override
  _CartState createState() => _CartState();
}

class _CartState extends State<Cart> {
  final String baseUrl = "http://localhost:3000";
  List<dynamic> products = [];

  @override
  void initState() {
    super.initState();
    getCart();
  }

  Future<void> getCart() async {
    final url = Uri.parse('$baseUrl/v1/orders/cart/?id=${widget.userId}');
    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {
        List jsonData = json.decode(response.body);
        setState(() {
          products = jsonData;
        });
      } else {
        print('Failed to load products: ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching products: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        toolbarHeight: 80,
        backgroundColor: Colors.red,
        title: Text(
          'Cart',
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
      body: Padding(
        padding: EdgeInsets.all(10),
        child: products.isEmpty
          ? Center(child: Text('Your cart is empty'))
          : ListView.builder(
          itemCount: products.length,
          itemBuilder: (context, index) {
            final product = products[index];
            int quantity = int.tryParse(product['quan_ordered']?.toString() ?? '0') ?? 0;
            print(quantity);
            return Container(
              padding: EdgeInsets.all(10),
              margin: EdgeInsets.only(bottom: 10),
              child: Row(
                children: [
                  Container(
                    width: 60,
                    height: 60,
                    decoration: BoxDecoration(
                      color: Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  SizedBox(width: 10),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            product['prod_name'],
                            style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                          ),
                          Text(
                            '₱${product['price'].toString()}',
                            style: TextStyle(fontSize: 16, color: Colors.red),
                          ),
                        ],
                      ),
                      SizedBox(height: 4),
                      Row(
                        children: [
                          IconButton(
                            onPressed: () {
                              if (quantity > 1) {
                                setState(() {
                                  products[index]['quan_ordered'] = quantity - 1;
                                });
                              }
                            },
                            icon: Icon(Icons.remove_circle, color: Colors.red),
                          ),
                          Text(
                            "${products[index]['quan_ordered']}", // Directly use updated value
                            style: TextStyle(fontSize: 16),
                          ),
                          IconButton(
                            onPressed: () {
                              setState(() {
                                products[index]['quan_ordered'] = quantity + 1;
                              });
                            },
                            icon: Icon(Icons.add_circle, color: Colors.red),
                          ),
                        ],
                      ),
                    ],
                  ),
                ],
              ),
            );
          }
        ),
      ),
      bottomNavigationBar: BottomNavBar(userId: widget.userId),
    );
  }
}
