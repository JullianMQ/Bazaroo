import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'cart.dart';
import 'nav/customer_nav.dart';

class ProductDetailScreen extends StatefulWidget {
  final String userId;
  final int prodId;  
  const ProductDetailScreen({Key? key, required this.userId, required this.prodId}) : super(key: key);

  @override
  _ProductDetailScreenState createState() => _ProductDetailScreenState();
}

class _ProductDetailScreenState extends State<ProductDetailScreen> {
  final String baseUrl = "https://bazaroo.onrender.com";
  bool isLoading = false; 
  int quantity = 1;
  Map<String, dynamic> product = {};

  @override
  void initState() {
    super.initState();
    fetchProducts();
  }

  Future<void> fetchProducts() async {
    setState(() {
      isLoading = true;
    });

    final url = Uri.parse('$baseUrl/v1/products/?id=${widget.prodId}');
    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {
        Map<String, dynamic> data = json.decode(response.body);
        setState(() {
          product = data;
          isLoading = false;
        });
      } else {
        setState(() {
          isLoading = false;
        });
      }
    } catch (e) {
      setState(() {
        isLoading = false;
      });
    }
  }

  Future<void> checkCart() async {
    final checkUrl = Uri.parse("$baseUrl/v1/orders/cart/?id=${widget.userId}");
    final response = await http.get(checkUrl);
    List cartItems = json.decode(response.body);
    bool found = false;

    for (var item in cartItems) {
      if (item['prod_id'] == widget.prodId) {
        int prevQuan = int.tryParse(item['quan_ordered']?.toString() ?? '0') ?? 0;
        int ordId = int.tryParse(item['ord_id']?.toString() ?? '0') ?? 0;
        updateCart(ordId, prevQuan + quantity);
        found = true;
        break;
      }
    }

    if (!found) {
      addToCart();
    }

    _showSuccessDialog("Successfully added to cart. Added $quantity ");
  }

  Future<void> updateCart(int ordId, int newQuan) async {
    final postUrl = Uri.parse("$baseUrl/v1/orderdetails/?ord_id=$ordId&prod_id=${widget.prodId}");
    await http.put(
      postUrl,
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({"quan_ordered": newQuan}),
    );
  }

  Future<void> addToCart() async {
    final addUrl = Uri.parse("$baseUrl/v1/orderdetails/addtocart/?id=${widget.userId}");
    await http.post(
      addUrl,
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({
        "cust_id": widget.userId,
        "prod_name": product['prod_name'] ?? "",
        "prod_id": widget.prodId,
        "price": product['buy_price'] ?? 0,
        "quan_ordered": quantity,
      }),
    );
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
                Navigator.pop(context);
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
        automaticallyImplyLeading: false,
        actions: [
          Container(
            margin: EdgeInsets.symmetric(horizontal: 25.0),
            child: Row(
              mainAxisSize: MainAxisSize.min, 
              children: [
                IconButton(
                  icon: Icon(Icons.shopping_cart, color: Colors.white, size: 25),
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(builder: (context) => Cart(userId: widget.userId)), 
                    );
                  },
                ),
                SizedBox(width: 10), 
                IconButton(
                  icon: Icon(Icons.chevron_right, color: Colors.white, size: 35),
                  onPressed: () {
                    Navigator.pop(context);
                  },
                ),
              ],
            ),
          ),
        ],
      ),
      body: isLoading
          ? Center(child: CircularProgressIndicator()) 
          : Padding(
              padding: EdgeInsets.all(20.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  SizedBox(height: 20),
                  AspectRatio(
                    aspectRatio: 1,
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(15),
                        image: DecorationImage(
                          image: NetworkImage("$baseUrl/v1/images/${product['prod_image']['String']}"),
                          fit: BoxFit.cover,
                        ),
                      ),
                    ),
                  ),
                  SizedBox(height: 10),
                  Align(
                    alignment: Alignment.centerLeft, 
                    child: Text(
                      product['prod_name'] ?? "Loading...",
                      style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                    ),
                  ),
                  SizedBox(height: 5),
                  Container( 
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          "₱${product['buy_price'] ?? 0}",
                          style: TextStyle(fontSize: 20, color: Colors.red),
                        ),
                        SizedBox(width: 20), 
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            IconButton(
                              onPressed: () {
                                if (quantity > 1) {
                                  setState(() {
                                    quantity--;
                                  });
                                }
                              },
                              icon: Icon(Icons.remove_circle, color: Colors.red),
                            ),
                            SizedBox(width: 10), 
                            Text(
                              "$quantity",
                              style: TextStyle(fontSize: 18),
                            ),
                            SizedBox(width: 10), 
                            IconButton(
                              onPressed: () {
                                setState(() {
                                  quantity++;
                                });
                              },
                              icon: Icon(Icons.add_circle, color: Colors.red),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                  SizedBox(height: 10),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        minimumSize: Size(double.infinity, 50),
                        backgroundColor: Colors.red,
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                        padding: EdgeInsets.symmetric(vertical: 12),
                      ),
                      onPressed: checkCart,
                      child: Text(
                        "Add to Cart",
                        style: TextStyle(color: Colors.white, fontSize: 16),
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
