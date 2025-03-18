import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'nav/customer_nav.dart';

class ProductDetailScreen extends StatefulWidget {
  final String userId;
  final int prodId;
  final String prodName;
  final int buyPrice;
  final String imgPath;
  
  const ProductDetailScreen({Key? key, required this.userId, required this.prodId, required this.prodName, required this.buyPrice,  required this.imgPath}) : super(key: key);

  @override
  _ProductDetailScreenState createState() => _ProductDetailScreenState();
}

class _ProductDetailScreenState extends State<ProductDetailScreen> {
  int quantity = 1;

  Future<void> addToCart() async {
    final url = Uri.parse("http://localhost:3000/v1/orderdetails/addtocart/?id=${widget.userId}");
    final response = await http.post(
      url,
      headers: {"Content-Type": "application/json"},
      body: jsonEncode({
        "cust_id": widget.userId,
        "prod_name": widget.prodName,
        "prod_id": widget.prodId,
        "price": widget.buyPrice,
        "quan_ordered": quantity,
      }),
    );

    if (response.statusCode == 200) {
      _showSuccessDialog("Added successfully! Quantity: $quantity");
    } else {
      print("Error: ${response.statusCode} - ${response.body}");
      _showSuccessDialog("Failed to add item. Try again.");
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
                    image: NetworkImage(widget.imgPath),
                    fit: BoxFit.cover,
                  ),
                ),
              ),
            ),
            SizedBox(height: 10),
            Align(
              alignment: Alignment.centerLeft, 
              child: Text(
                widget.prodName,
                style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
              ),
            ),
            SizedBox(height: 5),
            Container( 
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    "₱ ${widget.buyPrice}",
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
                onPressed: addToCart,
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
