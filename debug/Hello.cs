using System;

public class Hello {
    public static void PrintArray(string prefix, string[] words) {
        Console.WriteLine(prefix + ": " + String.Join(", ", words));
    }

    public static void Main(string[] args) {
        Console.WriteLine("Hello, World!");
        PrintArray("Before", args);
        Array.Sort(args);
        PrintArray("After", args);
    }
}
