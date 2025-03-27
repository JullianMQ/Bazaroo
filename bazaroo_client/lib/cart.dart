import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'nav/customer_nav.dart';
import 'productsPage.dart';

class Cart extends StatefulWidget {
  final String userId;
  const Cart({Key? key, required this.userId}) : super(key: key);

  @override
  _CartState createState() => _CartState();
}

class _CartState extends State<Cart> {
  final String baseUrl = "https://bazaroo.onrender.com";
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
  
  Future<void> updateQuantity(int ordId, int prodId, int quantity) async {
    final url = Uri.parse('$baseUrl/v1/orderdetails/?ord_id=$ordId&prod_id=$prodId');
    http.put(url, headers: {"Content-Type": "application/json"},
      body: jsonEncode({
        "quan_ordered": quantity,
      }),
    );
  }

  Future<void> prodCheckout() async {
    final url = Uri.parse('$baseUrl/v1/orders/checkout/?id=${widget.userId}');
    final response = await http.post(
      url,
      headers: {"Content-Type": "application/json"},
    );
    if (response.statusCode == 200) {
      _showSuccessDialog("Successfully checked out item.");
    } else {
      print("Error: ${response.statusCode} - ${response.body}");
      _showSuccessDialog("Failed to checkout item. Try again.");    
    }
  }
  void _showSuccessDialog(String message) {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text("Success", style: TextStyle(color: Colors.green)),
          content: Text(message),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pushReplacement(
                  context,
                  MaterialPageRoute(builder: (context) => Cart(userId: widget.userId)),
                );
              },
              child: const Text("OK"),
            ),
          ],
        );
      },
    );
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
      child: Column(
        children: [
          Expanded(
            child: Container(
              child: products.isEmpty
            ? Center(child: Text('Your cart is empty'))
            : SingleChildScrollView(
                child: Column(
                  children: [
                    ListView.builder(
                      shrinkWrap: true, 
                      physics: NeverScrollableScrollPhysics(), 
                      itemCount: products.length,
                      itemBuilder: (context, index) {
                        final product = products[index];
                        int quantity = int.tryParse(product['quan_ordered']?.toString() ?? '0') ?? 0;
                        int price = int.tryParse(product['price']?.toString() ?? '0') ?? 0;
                        int total = quantity * price;
                        return Container(
                          padding: EdgeInsets.all(10),
                          margin: EdgeInsets.only(bottom: 10),
                          child: Row(
                            children: [
                              Container(
                                width: 60,
                                height: 60,
                                decoration: BoxDecoration(
                                  borderRadius: BorderRadius.circular(8),
                                ),
                                child: ClipRRect(
                                  borderRadius: BorderRadius.circular(8),
                                  child: Image.network(
                                    baseUrl + product['prod_image']['String'],                                   
                                    fit: BoxFit.cover, 
                                    errorBuilder: (context, error, stackTrace) {
                                      return Container(
                                        color: Colors.grey.shade300,      
                                        child: Icon(Icons.broken_image, color: Colors.grey), 
                                      );
                                    },
                                  ),
                                ),
                              ),
                              SizedBox(width: 10),
                              Expanded(
                                child: Column(
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
                                          '₱$total',
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
                                              updateQuantity(products[index]['ord_id'], products[index]['prod_id'], quantity - 1);
                                            }
                                          },
                                          icon: Icon(Icons.remove_circle, color: Colors.red),
                                        ),
                                        SizedBox(width: 10), 
                                        Text(
                                          "$quantity",
                                          style: TextStyle(fontSize: 16),
                                        ),
                                        SizedBox(width: 10), 
                                        IconButton(
                                          onPressed: () {
                                            setState(() {
                                              products[index]['quan_ordered'] = quantity + 1;
                                            });
                                            updateQuantity(products[index]['ord_id'], products[index]['prod_id'], quantity + 1);
                                          },
                                          icon: Icon(Icons.add_circle, color: Colors.red),
                                        ),
                                      ],
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
                        );
                      },
                    ),
                  ],
                ),
              ),
            ),
          ),
          Container(
            child: products.isEmpty
            ? ElevatedButton(
              style: ElevatedButton.styleFrom(
                minimumSize: Size(double.infinity, 50),
                backgroundColor: Colors.red,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
                padding: EdgeInsets.symmetric(vertical: 12),
              ),
              onPressed: () => Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => ProductsPage(userId: widget.userId)),
              ),
              child: Text(
                'View Products',
                style: TextStyle(color: Colors.white, fontSize: 16),
              ),
            )
            :Padding(
              padding: EdgeInsets.only(top: 10),
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  minimumSize: Size(double.infinity, 50),
                  backgroundColor: Colors.red,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                  padding: EdgeInsets.symmetric(vertical: 12),
                ),
                onPressed: prodCheckout,
                child: Text(
                  'Checkout',
                  style: TextStyle(color: Colors.white, fontSize: 16),
                ),
              ),
            ),
          ),
        ],
      ),
    ),
      bottomNavigationBar: BottomNavBar(userId: widget.userId),
    );
  }
}
